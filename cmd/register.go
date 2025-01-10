package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/ini.v1"
)

type OpenIDConfig struct {
	AuthorizationEndpoint string `json:"authorization_endpoint" ini:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint" ini:"token_endpoint"`
	Issuer                string `json:"issuer" ini:"issuer"`
	ClientID              string `json:"client_id" ini:"client_id"`
	Scope                 string `json:"scope" ini:"scope"`
}

func init() {

	RegisterCommand(&Command{
		Name:        "register",
		Description: "DH CLI register ",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     registerHandler,
	})

}

func registerHandler(args []string, fs *flag.FlagSet) {
	if len(args) < 2 {
		log.Fatalf("Error: Name of configuration to register and authorization URL are required positional arguments.\nUsage: dhcli register <name> <url>")
	}

	sectionName := args[0]
	authUrl := args[1]
	iniPath := getIniPath()

	// Read or initialize ini file
	cfg, err := ini.Load(iniPath)
	if err != nil {
		cfg = ini.Empty()
	}

	// Fetch OpenID configuration and write to ini file
	openIDConfig := fetchOpenIDConfig("https://" + authUrl + "/.well-known/openid-configuration")
	cfg.Section(sectionName).ReflectFrom(&openIDConfig)
	err = cfg.SaveTo(iniPath)
	if err != nil {
		fmt.Printf("Failed to write ini file: %v", err)
		os.Exit(1)
	}
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

func getIniPath() string {
	iniPath, err := os.UserHomeDir()
	if err != nil {
		iniPath = "./"
	}
	iniPath += "/.cli.ini"

	return iniPath
}
