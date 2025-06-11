package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"dhcli/utils"
)

func Register(endpoint string, envName string) error {
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	cfg := utils.LoadIni(true)

	res, coreConfig := fetchCoreConfig(endpoint + ".well-known/configuration")

	if envName == "" || envName == "null" {
		envName = coreConfig.Name
		if envName == "" {
			return fmt.Errorf("environment name not provided and not found in core configuration")
		}
	}

	if cfg.HasSection(envName) {
		log.Printf("Section '%v' already exists, will be overwritten.\n", envName)
	}
	section := cfg.Section(envName)
	for _, k := range section.Keys() {
		section.DeleteKey(k.Name())
	}
	section.ReflectFrom(&coreConfig)

	openIDConfig := fetchOpenIDConfig(endpoint + ".well-known/openid-configuration")
	openIDConfig.ClientID = coreConfig.ClientID
	section.ReflectFrom(&openIDConfig)

	apiLevel := ""
	for k, v := range res {
		if !section.HasKey(k) && k != "dhcore_client_id" {
			section.NewKey(k, utils.ReflectValue(v))
		}
		if k == utils.ApiLevelKey {
			apiLevel = v.(string)
		}
	}

	apiLevelInt, err := strconv.Atoi(apiLevel)
	if err != nil {
		log.Println("WARNING: API level is not a valid integer.")
	} else if apiLevelInt < utils.MinApiLevel {
		log.Printf("WARNING: API level %v is below the CLI's minimum requirement %v.\n", apiLevelInt, utils.MinApiLevel)
	}

	section.NewKey(utils.UpdatedEnvKey, time.Now().Format(time.RFC3339))

	defaultSection := cfg.Section("DEFAULT")
	if !defaultSection.HasKey(utils.CurrentEnvironment) {
		defaultSection.NewKey(utils.CurrentEnvironment, envName)
	}

	addIniToGitignore()
	utils.SaveIni(cfg)

	log.Printf("'%v' registered.\n", envName)
	return nil
}

func fetchCoreConfig(configURL string) (map[string]interface{}, CoreConfig) {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Fatalf("Error fetching core config: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Core responded with status %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading core config response: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error parsing core config JSON: %v", err)
	}

	var config CoreConfig
	if err := json.Unmarshal(body, &config); err != nil {
		log.Fatalf("Error mapping core config: %v", err)
	}

	return data, config
}

func fetchOpenIDConfig(configURL string) OpenIDConfig {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Fatalf("Error fetching OpenID config: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("OpenID config error: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading OpenID config: %v", err)
	}

	var config OpenIDConfig
	if err := json.Unmarshal(body, &config); err != nil {
		log.Fatalf("Error parsing OpenID config JSON: %v", err)
	}

	return config
}

func addIniToGitignore() {
	const gitignorePath = "./.gitignore"
	f, err := os.OpenFile(gitignorePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Cannot open .gitignore: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == utils.IniName {
			return
		}
	}
	if _, err := f.WriteString(utils.IniName + "\n"); err != nil {
		log.Fatalf("Error writing to .gitignore: %v", err)
	}
}
