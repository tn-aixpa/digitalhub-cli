package service

import (
	"dhcli/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func RefreshHandler(env string) {
	// Read config from ini file
	cfg, section := utils.LoadIniConfig([]string{env})

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", section.Key("client_id").Value())
	data.Set("refresh_token", section.Key("refresh_token").Value())

	resp, err := http.Post(section.Key("token_endpoint").Value(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error refreshing token: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Token server error: %s\nBody: %s\n", resp.Status, string(body))
		os.Exit(1)
	}

	var responseJson map[string]interface{}
	json.Unmarshal(body, &responseJson)
	utils.UpdateKey(section, "access_token", responseJson["access_token"].(string))
	utils.UpdateKey(section, "refresh_token", responseJson["refresh_token"].(string))

	utils.SaveIni(cfg)
	log.Printf("Token refreshed.\n")
}
