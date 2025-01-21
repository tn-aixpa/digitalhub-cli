package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

func init() {

	RegisterCommand(&Command{
		Name:        "refresh",
		Description: "DH CLI refresh",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("e", "", "environment")
		},
		Handler: refreshHandler,
	})

}

func refreshHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	// Read config from ini file
	cfg, err := ini.Load(getIniPath())
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	fs.Parse(args)
	sectionName := fs.Lookup("e").Value.String()
	if sectionName == "" {
		defaultEnvironment := getDefaultEnvironment(cfg)
		if defaultEnvironment == "" {
			log.Fatalf("Error: environment flag (-e) was not passed and default environment is not specified in ini file.\nUsage: dhcli login -e <environment>")
		}
		sectionName = defaultEnvironment
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		log.Fatalf("Failed to read section '%s': %v.", sectionName, err)
	}

	openIDConfig := new(OpenIDConfig)
	section.MapTo(openIDConfig)

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", openIDConfig.ClientID)
	data.Set("refresh_token", openIDConfig.RefreshToken)

	resp, err := http.Post(openIDConfig.TokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("Error refreshing token: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Token server error: %s\nBody: %s", resp.Status, string(body))
	}

	var responseJson map[string]interface{}
	json.Unmarshal(body, &responseJson)
	openIDConfig.AccessToken = responseJson["access_token"].(string)
	openIDConfig.RefreshToken = responseJson["refresh_token"].(string)

	section.ReflectFrom(&openIDConfig)
	err = cfg.SaveTo(getIniPath())
	if err != nil {
		fmt.Printf("Failed to update ini file: %v", err)
		os.Exit(1)
	}
}
