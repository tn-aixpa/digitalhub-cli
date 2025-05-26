package cmd

import (
	"dhcli/utils"
	"encoding/json"
	"flag"
	"log"
	"os"

	"sigs.k8s.io/yaml"
)

func init() {
	RegisterCommand(&Command{
		Name:        "create",
		Description: "dhcli create [-e <environment> -p <project> -f <file> -n name --reset-id] <resource> <name>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("e", "", "environment")
			fs.String("p", "", "project")
			fs.String("f", "", "file")
			fs.String("n", "", "name")
			fs.Bool("reset-id", false, "reset ID, ignoring the one in the file")
		},
		Handler: createHandler,
	})
}

func createHandler(args []string, fs *flag.FlagSet) {
	fs.Parse(args)
	if len(fs.Args()) < 1 {
		log.Println("Error: resource type is required.")
		os.Exit(1)
	}
	resource := utils.TranslateEndpoint(fs.Args()[0])

	environment := fs.Lookup("e").Value.String()
	project := fs.Lookup("p").Value.String()
	filePath := fs.Lookup("f").Value.String()
	name := fs.Lookup("n").Value.String()
	resetId := fs.Lookup("reset-id").Value.String()

	// Validate parameters
	if resource != "projects" {
		if project == "" {
			log.Println("Project is mandatory when performing this operation on resources other than projects.")
			os.Exit(1)
		}
		if filePath == "" {
			log.Println("Input file not specified.")
			os.Exit(1)
		}
	} else if filePath == "" && name == "" {
		log.Println("Must provide either an input file or a name when creating a project.")
		os.Exit(1)
	}

	var jsonMap map[string]interface{}

	if filePath != "" {
		// Read file
		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read YAML file: %v\n", err)
			os.Exit(1)
		}

		// Convert YAML to JSON
		jsonBytes, err := yaml.YAMLToJSON(file)

		// Convert to map
		err = json.Unmarshal(jsonBytes, &jsonMap)
		if err != nil {
			log.Printf("Failed to parse after JSON conversion: %v\n", err)
			os.Exit(1)
		}

		// Alter fields
		delete(jsonMap, "user")

		if resource != "projects" {
			jsonMap["project"] = project
		}

		if resetId == "true" {
			delete(jsonMap, "id")
		}
	} else {
		jsonMap = map[string]interface{}{}
		jsonMap["name"] = name
	}

	// Marshal
	jsonBody, err := json.Marshal(jsonMap)
	if err != nil {
		log.Printf("Failed to marshal: %v\n", err)
		os.Exit(1)
	}

	// Request
	_, section := loadIniConfig([]string{environment})
	method := "POST"
	url := utils.BuildCoreUrl(section, project, resource, "", nil)
	req := utils.PrepareRequest(method, url, jsonBody, section.Key("access_token").String())
	utils.DoRequest(req)
	log.Println("Created successfully.")
}
