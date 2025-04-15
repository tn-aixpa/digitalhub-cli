package cmd

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"slices"
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
		Description: "dhcli login <environment>",
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
	fmt.Println("The following URL will be opened in your browser to authenticate:")
	fmt.Println(authURL)
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Press enter to continue...")
	_, err := buf.ReadBytes('\n')
	if err != nil {
		fmt.Printf("Error during confirmation: %v\n", err)
	}

	// Open the URL in the default browser
	err = openBrowser(authURL)
	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
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
			fmt.Printf("Error generating random string: %v\n", err)
			os.Exit(1)
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
			fmt.Printf("State mismatch: expected %s, got %s\n", generatedState, state)
			os.Exit(1)
		}

		if authCode == "" {
			http.Error(w, "Authorization code not received", http.StatusBadRequest)
			return
		}

		slog.Debug("Authorization code received correctly.", "Code", authCode, "State", state)

		tokenResponse := exchangeAuthCode(openIDConfig.TokenEndpoint, openIDConfig.ClientID, codeVerifier, authCode)
		if tokenResponse == nil {
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}

		// Give feedback to user
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<div style=\"margin: 24px 0px 0px 24px;\">")
		fmt.Fprintf(w, "<h1>Authorization successful</h1>")
		fmt.Fprintf(w, `<h3>You may now close this window.</h3>`)
		fmt.Fprintf(w, `<h3>Token response:</h3>`)
		fmt.Fprintf(w, `<div>`)
		fmt.Fprintf(w, "<button style=\"position: absolute;left: 774px;padding: 10px;opacity: 0.95;cursor: pointer;\" onclick=\"navigator.clipboard.writeText(document.getElementById('resp').innerHTML)\">Copy</button>")
		fmt.Fprintf(w, "<div id=\"resp\" style=\"width: 800px;font-family: courier;overflow: auto;height: 320px;overflow-wrap: break-word;\">%s</div>", tokenResponse)
		fmt.Fprintf(w, "</div>")
		fmt.Fprintf(w, "</div>")

		// Save response token
		slog.Debug("Token response received correctly.", "Response", string(tokenResponse))
		var responseJson map[string]interface{}
		json.Unmarshal(tokenResponse, &responseJson)
		for k, v := range responseJson {
			if !slices.Contains([]string{"client_id", "token_type", "id_token"}, k) {
				if !section.HasKey(k) {
					section.NewKey(k, utils.ReflectValue(v))
				} else {
					section.Key(k).SetValue(utils.ReflectValue(v))
				}
			}
		}
		openIDConfig.AccessToken = responseJson["access_token"].(string)
		refreshToken, ok := responseJson["refresh_token"]
		if ok {
			openIDConfig.RefreshToken = refreshToken.(string)
		}

		section.ReflectFrom(&openIDConfig)
		utils.SaveIni(cfg)
		fmt.Println("Login successful!")

		// Close cli immediately in a goroutine, this keeps the browser open but releases the command line tool
		go func() {
			os.Exit(0)
		}()
	})
	go func() {
		if err := http.ListenAndServe(":4000", nil); err != nil {
			slog.Error("Error starting server.", "Message", err)
			os.Exit(1)
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
	v.Set("code_challenge", codeChallenge)
	v.Set("code_challenge_method", "S256")
	v.Set("state", state)

	scopesString := strings.Join(openIDConfig.Scope[:], "%20")

	return fmt.Sprintf("%s?%s&scope=%s", openIDConfig.AuthorizationEndpoint, v.Encode(), scopesString)
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
		slog.Error("Error exchanging auth code for token.", "Message", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading token response.", "Message", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Token server error.", "Status", resp.Status, "Body", string(body))
		return nil
	}

	return body
}

// openBrowser tries to open the provided URL in the default web browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	// For Windows
	if runtime.GOOS == "windows" {
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
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
				fmt.Printf("Error while reading default environment: %v\n", err)
				os.Exit(1)
			}
			if defaultSection.HasKey("current_environment") {
				sectionName = defaultSection.Key("current_environment").String()
			}
		}

		if sectionName == "" {
			fmt.Println("Error: environment was not passed and default environment is not specified in ini file.")
			os.Exit(1)
		}
	} else {
		sectionName = args[0]
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		fmt.Printf("Failed to read section '%s': %v.\n", sectionName, err)
		os.Exit(1)
	}

	return cfg, section
}
