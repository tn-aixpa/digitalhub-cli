package cmd

import (
	"dhcli/utils"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
	"sigs.k8s.io/yaml"
)

func init() {
	RegisterCommand(&Command{
		Name:        "update",
		Description: "dhcli update [-e <environment> -p <project> -f <file>] <resource> <id>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("e", "", "environment")
			fs.String("p", "", "project")
			fs.String("f", "", "file")
		},
		Handler: updateHandler,
	})
}

func updateHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	fs.Parse(args)
	if len(fs.Args()) < 2 {
		fmt.Println("Error: resource type and id are required.")
		os.Exit(1)
	}
	resource := utils.TranslateEndpoint(fs.Args()[0])
	id := fs.Args()[1]

	environment := fs.Lookup("e").Value.String()
	project := fs.Lookup("p").Value.String()
	filePath := fs.Lookup("f").Value.String()

	if filePath == "" {
		fmt.Println("Input file not specified.")
		os.Exit(1)
	}

	if resource != "projects" && project == "" {
		fmt.Println("Project is mandatory when performing this operation on resources other than projects.")
		os.Exit(1)
	}
	// Read file
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read YAML file: %v\n", err)
		os.Exit(1)
	}

	// Convert YAML to JSON
	jsonBytes, err := yaml.YAMLToJSON(file)

	// Convert to map
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		fmt.Printf("Failed to parse after JSON conversion: %v\n", err)
		os.Exit(1)
	}

	// Alter fields
	if jsonMap["id"] != nil && jsonMap["id"] != id {
		fmt.Printf("Error: specified ID (%v) and ID found in file (%v) do not match. Are you sure you are trying to update the correct resource?\n", id, jsonMap["id"])
		os.Exit(1)
	}

	delete(jsonMap, "user")

	if resource != "projects" {
		jsonMap["project"] = project
	}

	// Marshal back
	jsonBody, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
		os.Exit(1)
	}

	// Request
	_, section := loadConfig([]string{environment})
	method := "PUT"
	url := utils.BuildCoreUrl(section, project, resource, id, nil)
	req := utils.PrepareRequest(method, url, jsonBody, section.Key("access_token").String())
	utils.DoRequest(req)
	fmt.Println("Updated successfully.")
}
