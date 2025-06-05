package cmd

import (
	"bufio"
	"dhcli/utils"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegisterCommand(&Command{
		Name:        "register",
		Description: "dhcli register [-e <environment>] <endpoint>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("e", "", "environment")
		},
		Handler: registerHandler,
	})
}

func registerHandler(args []string, fs *flag.FlagSet) {
	if len(args) < 1 {
		log.Println("Error: Endpoint is required.\nUsage: dhcli register [-e <environment name>] <endpoint>")
		os.Exit(1)
	}
	fs.Parse(args)
	environment := fs.Lookup("e").Value.String()
	endpoint := fs.Args()[0]
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	// Read or initialize ini file
	cfg := utils.LoadIni(true)

	//collect to map+struct
	res, coreConfig := fetchConfig(endpoint + ".well-known/configuration")
	if environment == "" || environment == "null" {
		environment = coreConfig.Name
		if environment == "" {
			log.Println("Failed to register: environment name not specified and not defined in core's configuration.")
			os.Exit(1)
		}
	}

	if cfg.HasSection(environment) {
		log.Printf("Section '%v' already exists, will be overwritten.\n", environment)
	}
	section := cfg.Section(environment)
	for _, k := range section.Keys() {
		section.DeleteKey(k.Name())
	}
	section.ReflectFrom(&coreConfig)

	// Fetch OpenID configuration
	openIDConfig := fetchOpenIDConfig(endpoint + ".well-known/openid-configuration")
	openIDConfig.ClientID = coreConfig.ClientID
	section.ReflectFrom(&openIDConfig)

	apiLevel := ""

	for k, v := range res {
		//add keys
		if !section.HasKey(k) && k != "dhcore_client_id" {
			section.NewKey(k, utils.ReflectValue(v))
		}

		if k == utils.ApiLevelKey {
			apiLevel = v.(string)
		}
	}

	apiLevelInt, err := strconv.Atoi(apiLevel)
	if err != nil {
		log.Println("WARNING: Registering an environment that may be incompatible with this version of the CLI: API level is not specified or cannot be read as integer.")
	} else if apiLevelInt < utils.MinApiLevel {
		log.Printf("WARNING: Registering an environment with an API level (%v) that does not meet the CLI's minimum requirement (%v). Some commands may not work correctly.\n", apiLevelInt, utils.MinApiLevel)
	}

	section.NewKey(utils.UpdatedEnvKey, time.Now().Format(time.RFC3339))

	// Check for default env
	defaultSection := cfg.Section("DEFAULT")
	if !defaultSection.HasKey(utils.CurrentEnvironment) {
		defaultSection.NewKey(utils.CurrentEnvironment, environment)
	}

	// gitignoreAddIniFile()
	utils.SaveIni(cfg)
	log.Printf("'%v' registered.\n", environment)
}

func fetchConfig(configURL string) (map[string]interface{}, utils.CoreConfig) {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Printf("Error fetching core configuration: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Core responded with error %v\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading core configuration response: %v\n", err)
		os.Exit(1)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("Error parsing core configuration: %v\n", err)
		os.Exit(1)
	}

	var config utils.CoreConfig
	if err := json.Unmarshal(body, &config); err != nil {
		log.Printf("Error parsing core configuration: %v\n", err)
		os.Exit(1)
	}

	return res, config
}

func fetchOpenIDConfig(configURL string) utils.OpenIDConfig {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Printf("Error fetching OpenID configuration: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Core responded with error %v\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading OpenID configuration response: %v\n", err)
		os.Exit(1)
	}

	var config utils.OpenIDConfig
	if err := json.Unmarshal(body, &config); err != nil {
		log.Printf("Error parsing OpenID configuration: %v\n", err)
		os.Exit(1)
	}

	return config
}

func gitignoreAddIniFile() {
	path := "./.gitignore"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Cannot open .gitignore file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == utils.IniName {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error while reading .gitignore file contents: %v\n", err)
		os.Exit(1)
	}

	if _, err = f.WriteString(utils.IniName); err != nil {
		log.Printf("Error while adding entry to .gitignore file: %v\n", err)
		os.Exit(1)
	}
}
