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
		Name:        "create",
		Description: "./dhcli create [-n <name> -p <project> -e <entity type>] <yaml file>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("n", "", "environment name")
			fs.String("p", "", "project")
			fs.String("e", "", "entity type")
		},
		Handler: createHandler,
	})
}

func createHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		log.Fatalf("Error: Path to YAML file is required.")
	}
	fs.Parse(args)
	name := fs.Lookup("n").Value.String()
	project := fs.Lookup("p").Value.String()
	entityType := fs.Lookup("e").Value.String()
	yamlFilePath := fs.Args()[0]

	if (project != "" && entityType == "") || (project == "" && entityType != "") {
		log.Fatalf("Cannot create entity unless both project and type are specified.")
	}

	_, section := loadConfig([]string{name})
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}
	jsonContents, err := yaml.YAMLToJSON(yamlFile)

	method := "POST"
	url := utils.BuildCoreUrl(section, method, project, entityType, "")
	req := utils.PrepareRequest(method, url, jsonContents, section.Key("access_token").String())
	utils.DoRequestAndPrintResponse(req)
}
