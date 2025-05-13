package cmd

import (
	"bytes"
	"dhcli/utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"sigs.k8s.io/yaml"
)

func init() {
	RegisterCommand(&Command{
		Name:        "get",
		Description: "dhcli get [-e <environment> -o <output format> -p <project> -n <name>] <resource> <id>",
		SetupFlags: func(fs *flag.FlagSet) {
			// CLI-specific
			fs.String("e", "", "environment")
			fs.String("o", "short", "output format")

			// API
			fs.String("p", "", "project")
			fs.String("n", "", "name")
		},
		Handler: getHandler,
	})
}

func getHandler(args []string, fs *flag.FlagSet) {
	fs.Parse(args)
	if len(fs.Args()) < 1 {
		log.Println("Error: resource type is required.")
		os.Exit(1)
	}
	resource := utils.TranslateEndpoint(fs.Args()[0])

	id := ""
	if len(fs.Args()) > 1 {
		id = fs.Args()[1]
	}

	environment := fs.Lookup("e").Value.String()
	format := utils.TranslateFormat(fs.Lookup("o").Value.String())
	project := fs.Lookup("p").Value.String()
	name := fs.Lookup("n").Value.String()

	if resource != "projects" && project == "" {
		log.Println("Project is mandatory when performing this operation on resources other than projects.")
		os.Exit(1)
	}

	params := map[string]string{}

	// If no ID is provided, search by name and latest version (result will be a list of 1)
	if id == "" {
		if name == "" {
			log.Println("You must specify id or name.")
			os.Exit(1)
		}
		params["name"] = name
		params["versions"] = "latest"
	}

	_, section := loadConfig([]string{environment})

	method := "GET"
	url := utils.BuildCoreUrl(section, project, resource, id, params)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())

	body, err := utils.DoRequest(req)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	printGet(format, id, body, args)
}

func getFirstIfList(m map[string]interface{}) map[string]interface{} {
	if content, ok := m["content"]; ok && reflect.ValueOf(content).Kind() == reflect.Slice {
		contentSlice := content.([]interface{})
		if len(contentSlice) >= 1 {
			return contentSlice[0].(map[string]interface{})
		}

		log.Println("Error: resource not found")
		os.Exit(1)
	}

	return m
}

func printGet(format string, id string, body []byte, args []string) {
	if format == "short" {
		printShortGet(body)
	} else if format == "json" {
		printJsonGet(id, body)
	} else if format == "yaml" {
		fmt.Printf("# Document generated with parameters: %v\n", strings.Join(args, " "))
		printYamlGet(id, body)
	}
}

func printShortGet(src []byte) {
	unmarshalled := map[string]interface{}{}
	json.Unmarshal([]byte(string(src)), &unmarshalled)
	m := getFirstIfList(unmarshalled)

	format := "%-12s %v\n"
	fmt.Printf(format, "Name:", m["name"])
	status := m["status"]
	if reflect.ValueOf(status).Kind() == reflect.Map {
		statusMap := status.(map[string]interface{})
		fmt.Printf(format, "State:", statusMap["state"])
	}
	fmt.Printf(format, "Kind:", m["kind"])
	fmt.Printf(format, "ID:", m["id"])
	fmt.Printf(format, "Key:", m["key"])

	metadata := m["metadata"]
	if reflect.ValueOf(metadata).Kind() == reflect.Map {
		metadataMap := metadata.(map[string]interface{})
		fmt.Printf(format, "Created on:", metadataMap["created"])
		fmt.Printf(format, "Created by:", metadataMap["created_by"])
		fmt.Printf(format, "Updated on:", metadataMap["updated"])
		fmt.Printf(format, "Updated by:", metadataMap["updated_by"])
	}
}

func printJsonGet(id string, src []byte) {
	resource := src

	// if ID is empty, result is a list, must extract first
	if id == "" {
		m := map[string]interface{}{}
		json.Unmarshal([]byte(string(src)), &m)
		first := getFirstIfList(m)
		remarshal, err := json.Marshal(first)
		if err != nil {
			log.Printf("Error while re-marshalling into JSON after extracting first element of list: %v\n", err)
		}
		resource = remarshal
	}

	var j bytes.Buffer
	err := json.Indent(&j, resource, "", "    ")
	if err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(j.Bytes()))
}

func printYamlGet(id string, src []byte) {
	resource := src

	// if ID is empty, result is a list, must extract first
	if id == "" {
		m := map[string]interface{}{}
		json.Unmarshal([]byte(string(src)), &m)
		first := getFirstIfList(m)
		remarshal, err := yaml.Marshal(first)
		if err != nil {
			log.Printf("Error while re-marshalling into YAML after extracting first element of list: %v\n", err)
		}
		resource = remarshal
	} else {
		converted, err := yaml.JSONToYAML(src)
		if err != nil {
			log.Printf("Error converting JSON into YAML: %v\n", err)
			os.Exit(1)
		}
		resource = converted
	}

	fmt.Println(string(resource))
}
