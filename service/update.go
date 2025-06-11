package service

import (
	"dhcli/utils"
	"encoding/json"
	"log"
	"os"
	"sigs.k8s.io/yaml"
)

func UpdateHandler(env, project, filePath, resource, id string, originalArgs []string) error {

	validResource := utils.TranslateEndpoint(resource)

	// Load environment and check API level requirements
	cfg, section := utils.LoadIniConfig([]string{env})
	utils.CheckUpdateEnvironment(cfg, section)
	utils.CheckApiLevel(section, utils.UpdateMin, utils.UpdateMax)

	if filePath == "" {
		log.Println("Input file not specified.")
		os.Exit(1)
	}

	if validResource != "projects" && project == "" {
		log.Println("Project is mandatory when performing this operation on resources other than projects.")
		os.Exit(1)
	}
	// Read a file
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read YAML file: %v\n", err)
		os.Exit(1)
	}

	// Convert YAML to JSON
	jsonBytes, err := yaml.YAMLToJSON(file)

	// Convert to map
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		log.Printf("Failed to parse after JSON conversion: %v\n", err)
		os.Exit(1)
	}

	// Alter fields
	if jsonMap["id"] != nil && jsonMap["id"] != id {
		log.Printf("Error: specified ID (%v) and ID found in file (%v) do not match. Are you sure you are trying to update the correct validResource?\n", id, jsonMap["id"])
		os.Exit(1)
	}

	delete(jsonMap, "user")

	if validResource != "projects" {
		jsonMap["project"] = project
	}

	// Marshal back
	jsonBody, err := json.Marshal(jsonMap)
	if err != nil {
		log.Printf("Failed to marshal: %v\n", err)
		os.Exit(1)
	}

	// Request
	method := "PUT"
	url := utils.BuildCoreUrl(section, project, validResource, id, nil)
	req := utils.PrepareRequest(method, url, jsonBody, section.Key("access_token").String())
	utils.DoRequest(req)
	log.Println("Updated successfully.")
	return nil
}
