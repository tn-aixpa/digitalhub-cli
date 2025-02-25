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

	"gopkg.in/ini.v1"

	"dhcli/utils"
)

var (
	redirectURI    = "http://localhost:4000/callback"
	generatedState string
)

func init() {
	RegisterCommand(&Command{
		Name:        "login",
		Description: "./dhcli login <environment>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     loginHandler,
	})
}

func loginHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	// Read config from ini file
	cfg, section := loadConfig(args)

	// Generate PKCE values
	codeVerifier, codeChallenge := generatePKCE()

	// Generate state value
	generatedState = generateRandomString(32)

	// Start local server to capture the authorization code
	startAuthCodeServer(cfg, section, codeVerifier)

	// Build and display the authorization URL
	authURL := buildAuthURL(section, codeChallenge, generatedState)
	fmt.Println("The following URL should open in your browser to authenticate:")
	fmt.Println(authURL)

	// Open the URL in the default browser
	err := openBrowser(authURL)
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}

	// Block the program to wait for user interaction
	select {}
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

func startAuthCodeServer(cfg *ini.File, section *ini.Section, codeVerifier string) {
	openIDConfig := new(OpenIDConfig)
	section.MapTo(openIDConfig)

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

		tokenResponse := exchangeAuthCode(openIDConfig.TokenEndpoint, openIDConfig.ClientID, codeVerifier, authCode)
		if tokenResponse == nil {
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}

		// Give feedback to user
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Authorization Successful</h1>")
		fmt.Fprintf(w, `<h2>Token response is:</h2>`)
		fmt.Fprintf(w, "<pre>%s</pre>", tokenResponse)
		fmt.Fprintf(w, `<h2>You may now close this window.</h2>`)

		// Save response token
		log.Println("Token Response:", string(tokenResponse))
		var responseJson map[string]interface{}
		json.Unmarshal(tokenResponse, &responseJson)
		openIDConfig.AccessToken = responseJson["access_token"].(string)
		refreshToken, ok := responseJson["refresh_token"]
		if ok {
			openIDConfig.RefreshToken = refreshToken.(string)
		}

		section.ReflectFrom(&openIDConfig)
		utils.SaveIni(cfg)

		// Close cli immediately in a goroutine, this keeps the browser open but releases the command line tool
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

func buildAuthURL(section *ini.Section, codeChallenge, state string) string {
	openIDConfig := new(OpenIDConfig)
	section.MapTo(openIDConfig)

	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", openIDConfig.ClientID)
	v.Set("redirect_uri", redirectURI)
	v.Set("scope", openIDConfig.Scope)
	v.Set("code_challenge", codeChallenge)
	v.Set("code_challenge_method", "S256")
	v.Set("state", state)

	return fmt.Sprintf("%s?%s", openIDConfig.AuthorizationEndpoint, v.Encode())
}

func exchangeAuthCode(tokenEndpoint, clientID, codeVerifier, authCode string) []byte {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("code_verifier", codeVerifier)
	data.Set("code", authCode)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error exchanging auth code for token: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading token response: %v", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Token server error: %s\nBody: %s", resp.Status, string(body))
		return nil
	}

	return body
}

// openBrowser tries to open the provided URL in the default web browser
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

func loadConfig(args []string) (*ini.File, *ini.Section) {
	cfg := utils.LoadIni(false)

	sectionName := ""

	if len(args) == 0 || args[0] == "" {
		if cfg.HasSection("DEFAULT") {
			defaultSection, err := cfg.GetSection("DEFAULT")
			if err != nil {
				log.Fatalf("Error while reading default environment: %v", err)
			}
			if defaultSection.HasKey("current_environment") {
				sectionName = defaultSection.Key("current_environment").String()
			}
		}

		if sectionName == "" {
			log.Fatalf("Error: environment was not passed and default environment is not specified in ini file.")
		}
	} else {
		sectionName = args[0]
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		log.Fatalf("Failed to read section '%s': %v.", sectionName, err)
	}

	return cfg, section
}
