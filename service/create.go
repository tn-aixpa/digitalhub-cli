package service

import (
	"dhcli/utils"
	"encoding/json"
	"log"
	"os"
	"sigs.k8s.io/yaml"
)

func CreateHandler(env, project string, name string, filePath string, resetId bool, resource, id string, originalArgs []string) error {

	endpoint := utils.TranslateEndpoint(resource)

	// Load environment and check API level requirements
	cfg, section := utils.LoadIniConfig([]string{env})
	utils.CheckUpdateEnvironment(cfg, section)
	utils.CheckApiLevel(section, utils.CreateMin, utils.CreateMax)

	// Validate parameters
	if endpoint != "projects" {
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

		if endpoint != "projects" {
			jsonMap["project"] = project
		}

		if resetId == true {
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
	method := "POST"
	url := utils.BuildCoreUrl(section, project, endpoint, "", nil)
	req := utils.PrepareRequest(method, url, jsonBody, section.Key("access_token").String())
	_, err = utils.DoRequest(req)
	if err != nil {
		return err
	}

	log.Println("Created successfully.")
	return nil
}
