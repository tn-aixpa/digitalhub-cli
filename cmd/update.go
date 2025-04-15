package cmd

import (
	"dhcli/utils"
	"flag"
	"log"
	"os"

	"gopkg.in/ini.v1"
	"sigs.k8s.io/yaml"
)

func init() {
	RegisterCommand(&Command{
		Name:        "update",
		Description: "./dhcli update [-n <name> -e <entity type> -i <id>] <project> <yaml file>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("n", "", "environment name")
			fs.String("e", "", "entity type")
			fs.String("i", "", "id")
		},
		Handler: updateHandler,
	})
}

func updateHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 2 {
		log.Fatalf("Error: Project and path to YAML file are required.")
	}
	fs.Parse(args)
	name := fs.Lookup("n").Value.String()
	project := fs.Args()[0]
	entityType := fs.Lookup("e").Value.String()
	id := fs.Lookup("i").Value.String()
	yamlFilePath := fs.Args()[1]

	if entityType != "" && id == "" {
		log.Fatalf("Entity type specified, but ID missing.")
	} else if entityType == "" && id != "" {
		log.Fatalf("ID specified, but entity type missing.")
	}

	_, section := loadConfig([]string{name})
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}
	jsonContents, err := yaml.YAMLToJSON(yamlFile)

	method := "PUT"
	url := utils.BuildCoreUrl(section, method, project, entityType, id)
	req := utils.PrepareRequest(method, url, jsonContents, section.Key("access_token").String())
	utils.DoRequestAndPrintResponse(req)
}
