// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"dhcli/utils"
)

const redirectURI = "http://localhost:4000/callback"

var generatedState string

// LoginHandler esegue il flusso PKCE per autenticarsi
func LoginHandler(env string) error {
	cfg, section := loadIniCfg(env)

	utils.CheckUpdateEnvironment(cfg, section)
	utils.CheckApiLevel(section, utils.LoginMin, utils.LoginMax)

	cv, cc := generatePKCE()
	generatedState = randomString(32)

	startAuthCodeServer(cfg, section, cv)

	authURL := buildAuthURL(section, cc, generatedState)
	// Log più leggibile in terminale
	fmt.Println("──────────────────────────────────────────────")
	fmt.Println("🔐  LoginHandler: Visit this URL to authenticate:")
	fmt.Println("──────────────────────────────────────────────")
	fmt.Println(authURL)
	fmt.Println("──────────────────────────────────────────────")
	fmt.Print("Press Enter to open browser… ")

	_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		fmt.Println("Error while authenticating...")
		return err
	}

	if err := openBrowser(authURL); err != nil {
		log.Printf("Error opening browser: %v", err)
	}

	select {} // blocca finché il server non chiude l'app
}

func loadIniCfg(env string) (*ini.File, *ini.Section) {
	return utils.LoadIniConfig([]string{env})
}

func generatePKCE() (verifier, challenge string) {
	const cs = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	verifier = randomStringCharset(64, cs)
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return
}

func randomString(n int) string {
	return randomStringCharset(n, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
}

func randomStringCharset(n int, cs string) string {
	b := make([]byte, n)
	for i := range b {
		_, _ = rand.Read(b[i : i+1])
		b[i] = cs[int(b[i])%len(cs)]
	}
	return string(b)
}

func startAuthCodeServer(cfg *ini.File, section *ini.Section, verifier string) {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		authCode := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if state != generatedState {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			log.Fatalf("State mismatch: got %q", state)
		}
		if authCode == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		tkn := exchangeAuthCode(
			section.Key("token_endpoint").String(),
			section.Key("client_id").String(),
			verifier,
			authCode,
		)
		if tkn == nil {
			http.Error(w, "Failed token exchange", http.StatusInternalServerError)
			return
		}

		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, tkn, "", "  "); err != nil {
			prettyJSON.Write(tkn) // fallback testo semplice se errore
		}

		fmt.Fprintln(w, "<h1>LoginHandler successful</h1>")
		fmt.Fprintf(w, "<pre style=\"background:#f6f8fa;border:1px solid #ccc;padding:16px;width:800px;overflow:auto;\">%s</pre>", prettyJSON.String())

		var m map[string]interface{}
		json.Unmarshal(tkn, &m)
		for k, v := range m {
			if !slices.Contains([]string{"client_id", "token_type", "id_token"}, k) {
				utils.UpdateKey(section, k, fmt.Sprint(v))
			}
		}
		utils.UpdateKey(section, "access_token", fmt.Sprint(m["access_token"]))
		if rt, ok := m["refresh_token"]; ok {
			utils.UpdateKey(section, "refresh_token", fmt.Sprint(rt))
		}
		utils.SaveIni(cfg)

		log.Println("LoginHandler successful.")
		go os.Exit(0)
	})
	go http.ListenAndServe(":4000", nil)
}

func exchangeAuthCode(tokenURL, clientID, verifier, code string) []byte {
	v := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"code_verifier": {verifier},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}
	resp, err := http.PostForm(tokenURL, v)
	if err != nil {
		log.Printf("Token request error: %v", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Token error %s: %s", resp.Status, body)
		return nil
	}
	tkn, _ := io.ReadAll(resp.Body)
	return tkn
}

func buildAuthURL(section *ini.Section, chal, state string) string {
	v := url.Values{
		"response_type":         {"code"},
		"client_id":             {section.Key("client_id").String()},
		"redirect_uri":          {redirectURI},
		"code_challenge":        {chal},
		"code_challenge_method": {"S256"},
		"state":                 {state},
	}
	scope := strings.ReplaceAll(section.Key("scopes_supported").String(), ",", "%20")
	return section.Key("authorization_endpoint").String() + "?" + v.Encode() + "&scope=" + scope
}

func openBrowser(u string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", u)
	case "darwin":
		cmd = exec.Command("open", u)
	default:
		cmd = exec.Command("xdg-open", u)
	}
	return cmd.Start()
}
