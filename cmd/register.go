// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bufio"
	"dhcli/utils"
	"flag"
	"log"
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

	cfg := utils.LoadIni(true)

	// Configuration
	config, err := utils.FetchConfig(endpoint + ".well-known/configuration")
	if err != nil {
		log.Printf("Error while fetching configuration: %v\n", err)
		os.Exit(1)
	}
	if environment == "" || environment == "null" {
		environment = utils.GetStringValue(config, "dhcore_name")
		if environment == "" {
			log.Println("Failed to register: environment name not specified and not defined in core's configuration.")
			os.Exit(1)
		}
	}

	// If section already exists, all existing keys have to be removed
	if cfg.HasSection(environment) {
		log.Printf("Section '%v' already exists, will be overwritten.\n", environment)
	}
	section := cfg.Section(environment)
	for _, k := range section.Keys() {
		section.DeleteKey(k.Name())
	}

	// Copy keys and values
	for k, v := range config {
		newKey := k
		if newKey == utils.ClientIdKey {
			newKey = "client_id"
		}
		section.NewKey(newKey, utils.ReflectValue(v))
	}

	// Check API level compatibility
	apiLevel := utils.GetStringValue(config, utils.ApiLevelKey)
	apiLevelInt, err := strconv.Atoi(apiLevel)
	if err != nil {
		log.Println("WARNING: Registering an environment that may be incompatible with this version of the CLI: API level is not specified or cannot be read as integer.")
	} else if apiLevelInt < utils.MinApiLevel {
		log.Printf("WARNING: Registering an environment with an API level (%v) that does not meet the CLI's minimum requirement (%v). Some commands may not work correctly.\n", apiLevelInt, utils.MinApiLevel)
	}

	// OpenID configuration
	openIdConfig, err := utils.FetchConfig(endpoint + ".well-known/openid-configuration")
	if err != nil {
		log.Printf("Error while fetching OpenID configuration: %v\n", err)
		os.Exit(1)
	}
	for _, k := range utils.OpenIdFields {
		var v interface{} = ""
		if val, ok := openIdConfig[k]; ok {
			v = val
		}
		section.NewKey(k, utils.ReflectValue(v))
	}

	// Timestamp
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
