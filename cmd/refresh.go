package cmd

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/ini.v1"

	"dhcli/utils"
)

func init() {

	RegisterCommand(&Command{
		Name:        "refresh",
		Description: "./dhcli refresh <environment>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     refreshHandler,
	})

}

func refreshHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	// Read config from ini file
	cfg, section := loadConfig(args)
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
	utils.SaveIni(cfg)
}
