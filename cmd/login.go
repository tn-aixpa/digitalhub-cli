package cmd

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	redirectURI    = "http://localhost:4000/callback"
	clientID       = "c_5327245aac20400897a20ed4d0deef86"
	generatedState string
)

type OpenIDConfig struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	Issuer                string `json:"issuer"`
	ClientID              string `json:"client_id"`
	Scope                 string `json:"scope"`
}

func init() {
	RegisterCommand(&Command{
		Name:        "login",
		Description: "DH CLI login",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     loginHandler,
	})
}

func loginHandler(args []string, fs *flag.FlagSet) {
	if len(args) < 1 {
		log.Fatalf("Error: URL is required as a positional argument.\nUsage: dhcli login <url>")
	}

	authUrl := args[0]
	fmt.Printf("Authorize with: %s\n", authUrl)

	// Step 1: Fetch OpenID configuration
	openIDConfig := fetchOpenIDConfig("https://" + authUrl + "/.well-known/openid-configuration")

	// Step 2: Generate PKCE values
	codeVerifier, codeChallenge := generatePKCE()

	// Step 3: Generate state value
	generatedState = generateRandomString(32)

	// Step 4: Start local server to capture the authorization code
	startAuthCodeServer(openIDConfig, codeVerifier)

	// Step 5: Build and display the authorization URL
	authURL := buildAuthURL(openIDConfig.AuthorizationEndpoint, clientID, openIDConfig.Scope, codeChallenge, generatedState)
	fmt.Println("Open the following URL in canyour browser to authenticate:")
	fmt.Println(authURL)

	// Open the URL in the default browser
	err := openBrowser(authURL)
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}

	// Block the program to wait for user interaction
	select {}
}

func fetchOpenIDConfig(configURL string) OpenIDConfig {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Fatalf("Error fetching OpenID configuration: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading OpenID configuration response: %v", err)
	}

	var config OpenIDConfig
	if err := json.Unmarshal(body, &config); err != nil {
		log.Fatalf("Error parsing OpenID configuration: %v", err)
	}

	return config
}

func generatePKCE() (string, string) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	const length = 64

	codeVerifier := generateRandomStringWithCharset(length, charset)
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return codeVerifier, codeChallenge
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	return generateRandomStringWithCharset(length, charset)
}

func generateRandomStringWithCharset(length int, charset string) string {
	result := make([]byte, length)
	for i := range result {
		randomByte := make([]byte, 1)
		if _, err := rand.Read(randomByte); err != nil {
			log.Fatalf("Error generating random string: %v", err)
		}
		result[i] = charset[randomByte[0]%byte(len(charset))]
	}
	return string(result)
}

func startAuthCodeServer(config OpenIDConfig, codeVerifier string) {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		authCode := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if state != generatedState {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			log.Fatalf("State mismatch: expected %s, got %s", generatedState, state)
		}

		if authCode == "" {
			http.Error(w, "Authorization code not received", http.StatusBadRequest)
			return
		}

		log.Printf("Authorization Code: %s, State: %s\n", authCode, state)

		tokenResponse := exchangeAuthCode(config.TokenEndpoint, clientID, codeVerifier, authCode)
		if tokenResponse == "" {
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}

		// Give feedback to user and close the window after 5 seconds..
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Authorization Successful</h1>")
		fmt.Fprintf(w, "<pre>%s</pre>", tokenResponse)
		fmt.Fprintf(w, `
			<h2>This page will automatically close in <span id="countdown">5</span> seconds.</h2>
			<script type="text/javascript">
				var countdown = 5; 
				var countdownElement = document.getElementById('countdown');
				
				// Update the countdown every second
				var interval = setInterval(function() {
					countdown--; 
					countdownElement.innerHTML = countdown; // update the countdown display
					
					if (countdown <= 0) {
						clearInterval(interval); 
						window.close(); 
					}
				}, 1000);
			</script>
		`)

		log.Println("Token Response:", tokenResponse)

		// Close cli immediately in a goroutine, this keep open the browser but release the command line tool
		go func() {
			os.Exit(0)
		}()
	})
	go func() {
		if err := http.ListenAndServe(":4000", nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
}

func buildAuthURL(authEndpoint, clientID, scope, codeChallenge, state string) string {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", clientID)
	v.Set("redirect_uri", redirectURI)
	v.Set("scope", scope)
	v.Set("code_challenge", codeChallenge)
	v.Set("code_challenge_method", "S256")
	v.Set("state", state)

	return fmt.Sprintf("%s?%s", authEndpoint, v.Encode())
}

func exchangeAuthCode(tokenEndpoint, clientID, codeVerifier, authCode string) string {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("code_verifier", codeVerifier)
	data.Set("code", authCode)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error exchanging auth code for token: %v", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading token response: %v", err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Token server error: %s\nBody: %s", resp.Status, string(body))
		return ""
	}

	return string(body)
}

// openBrowser tries to open the provided URL in the default web browser.
func openBrowser(url string) error {
	var cmd *exec.Cmd

	// For Windows
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", url)
	} else if runtime.GOOS == "darwin" { // For macOS
		cmd = exec.Command("open", url)
	} else { // For Linux and others
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
