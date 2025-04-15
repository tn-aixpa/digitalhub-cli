package cmd

import (
	"dhcli/utils"
	"flag"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
	"sigs.k8s.io/yaml"
)

func init() {
	RegisterCommand(&Command{
		Name:        "create",
		Description: "dhcli create [-n <name> -p <project> -e <entity type>] <yaml file>",
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
		fmt.Println("Error: Path to YAML file is required.")
		os.Exit(1)
	}
	fs.Parse(args)
	name := fs.Lookup("n").Value.String()
	project := fs.Lookup("p").Value.String()
	entityType := fs.Lookup("e").Value.String()
	yamlFilePath := fs.Args()[0]

	if (project != "" && entityType == "") || (project == "" && entityType != "") {
		fmt.Println("Cannot create entity unless both project and type are specified.")
		os.Exit(1)
	}

	_, section := loadConfig([]string{name})
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		fmt.Printf("Failed to read YAML file: %v\n", err)
		os.Exit(1)
	}
	jsonContents, err := yaml.YAMLToJSON(yamlFile)

	method := "POST"
	url := utils.BuildCoreUrl(section, method, project, entityType, "")
	req := utils.PrepareRequest(method, url, jsonContents, section.Key("access_token").String())
	utils.DoRequestAndPrintResponse(req)
}
