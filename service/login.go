// service/login.go
package service

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

const redirectURI = "http://localhost:4000/callback"

var generatedState string

type OpenIDConfig struct {
	ClientID              string   `ini:"client_id"`
	AuthorizationEndpoint string   `ini:"authorization_endpoint"`
	TokenEndpoint         string   `ini:"token_endpoint"`
	Scope                 []string `delim:" " ini:"scope"`
	AccessToken           string
	RefreshToken          string
}

func Login(environment string) error {
	cfg := utils.LoadIni(false)

	sectionName := environment
	if environment == "" {
		defaultSection, _ := cfg.GetSection("DEFAULT")
		if defaultSection.HasKey("current_environment") {
			sectionName = defaultSection.Key("current_environment").String()
		}
		if sectionName == "" {
			return fmt.Errorf("environment not passed and no default found in ini")
		}
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return fmt.Errorf("failed to read section '%s': %w", sectionName, err)
	}

	utils.CheckApiLevel(section, utils.LoginMin, utils.LoginMax)

	codeVerifier, codeChallenge := generatePKCE()
	generatedState = generateRandomString(32)

	startAuthCodeServer(cfg, section, codeVerifier)

	authURL := buildAuthURL(section, codeChallenge, generatedState)
	log.Println("The following URL will be opened in your browser to authenticate:")
	log.Println(authURL)
	log.Println("Press enter to continue...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')

	if err := openBrowser(authURL); err != nil {
		log.Printf("Error opening browser: %v\n", err)
	}

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
		_, _ = rand.Read(randomByte)
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

		tokenResponse := exchangeAuthCode(openIDConfig.TokenEndpoint, openIDConfig.ClientID, codeVerifier, authCode)
		if tokenResponse == nil {
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}

		// Feedback HTML
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div style="margin: 24px;"><h1>Authorization successful</h1><h3>You may now close this window.</h3><h3>Token response:</h3>`)
		fmt.Fprintf(w, `<div><button onclick="navigator.clipboard.writeText(document.getElementById('resp').innerHTML)">Copy</button>`)
		fmt.Fprintf(w, `<div id="resp" style="font-family: courier; overflow: auto; height: 320px;">%s</div></div></div>`, tokenResponse)

		var responseJson map[string]interface{}
		json.Unmarshal(tokenResponse, &responseJson)
		for k, v := range responseJson {
			if !slices.Contains([]string{"client_id", "token_type", "id_token"}, k) {
				val := utils.ReflectValue(v)
				if !section.HasKey(k) {
					section.NewKey(k, val)
				} else {
					section.Key(k).SetValue(val)
				}
			}
		}

		openIDConfig.AccessToken = responseJson["access_token"].(string)
		if rtk, ok := responseJson["refresh_token"]; ok {
			openIDConfig.RefreshToken = rtk.(string)
		}

		section.ReflectFrom(openIDConfig)
		utils.SaveIni(cfg)

		log.Println("Login successful!")
		go func() {
			os.Exit(0)
		}()
	})

	go func() {
		if err := http.ListenAndServe(":4000", nil); err != nil {
			slog.Error("Error starting server", "Message", err)
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

	scopesString := strings.Join(openIDConfig.Scope, "%20")
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
		slog.Error("Error exchanging code for token", "Message", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading token response", "Message", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Token exchange failed", "Status", resp.Status, "Body", string(body))
		return nil
	}

	return body
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
