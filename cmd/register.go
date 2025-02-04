package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/ini.v1"

	"dhcli/utils"
)

type OpenIDConfig struct {
	AuthorizationEndpoint string `json:"authorization_endpoint" ini:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint" ini:"token_endpoint"`
	Issuer                string `json:"issuer" ini:"issuer"`
	ClientID              string `json:"client_id" ini:"client_id"`
	Scope                 string `json:"scope" ini:"scope"`
	AccessToken           string `json:"access_token" ini:"access_token"`
	RefreshToken          string `json:"refresh_token" ini:"refresh_token"`
}

func init() {
	RegisterCommand(&Command{
		Name:        "register",
		Description: "./dhcli register [-s <scope>] <environment> <authorization_provider> <client_id>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("s", "", "scope")
		},
		Handler: registerHandler,
	})
}

func registerHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 3 {
		log.Fatalf("Error: The following positional parameters are required: environment name, authorization URL, client ID.\nUsage: ./dhcli register [-s <scope>] <environment> <authorization_provider> <client_id>")
	}
	fs.Parse(args)
	scope := fs.Lookup("s").Value.String()

	sectionName := fs.Args()[0]
	authUrl := fs.Args()[1]
	clientId := fs.Args()[2]

	// Read or initialize ini file
	cfg := utils.LoadIni(true)

	// Fetch OpenID configuration and write to ini file
	openIDConfig := fetchOpenIDConfig("https://" + authUrl + "/.well-known/openid-configuration")
	openIDConfig.ClientID = clientId
	openIDConfig.Scope = scope

	cfg.Section(sectionName).ReflectFrom(&openIDConfig)
	gitignoreAddIniFile()
	utils.SaveIni(cfg)
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

func gitignoreAddIniFile() {
	path := "./.gitignore"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Cannot open .gitignore file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == utils.IniName {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading .gitignore file contents: %v", err)
	}

	if _, err = f.WriteString(utils.IniName); err != nil {
		log.Fatalf("Error while adding entry to .gitignore file: %v", err)
	}
}
