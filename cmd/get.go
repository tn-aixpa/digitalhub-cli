package cmd

import (
	"bytes"
	"dhcli/utils"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/ini.v1"
	"sigs.k8s.io/yaml"
)

func init() {
	RegisterCommand(&Command{
		Name:        "get",
		Description: "dhcli get [-e <environment> -o <output format> -p <project>] <resource> <id>",
		SetupFlags: func(fs *flag.FlagSet) {
			// CLI-specific
			fs.String("e", "", "environment")
			fs.String("o", "short", "output format")

			// API
			fs.String("p", "", "project")
		},
		Handler: getHandler,
	})
}

func getHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	fs.Parse(args)
	if len(fs.Args()) < 2 {
		fmt.Println("Error: resource type and id are required.")
		os.Exit(1)
	}
	resource := utils.TranslateEndpoint(fs.Args()[0])
	id := fs.Args()[1]

	environment := fs.Lookup("e").Value.String()
	outputFormat := fs.Lookup("o").Value.String()
	project := fs.Lookup("p").Value.String()

	if resource != "projects" && project == "" {
		fmt.Println("Project is mandatory when performing this operation on resources other than projects.")
		os.Exit(1)
	}

	_, section := loadConfig([]string{environment})

	method := "GET"
	url := utils.BuildCoreUrl(section, project, resource, id, nil)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())

	body, err := utils.DoRequest(req)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	format := utils.TranslateFormat(outputFormat)
	if format == "short" {
		mapResp := map[string]interface{}{}
		json.Unmarshal([]byte(string(body)), &mapResp)
		printShortGet(mapResp)
	} else if format == "json" {
		printJsonGet(body)
	} else if format == "yaml" {
		printYamlGet(body)
	}
}

func printShortGet(m map[string]interface{}) {
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

func printJsonGet(src []byte) {
	var j bytes.Buffer
	err := json.Indent(&j, src, "", "    ")
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(j.Bytes()))
}

func printYamlGet(src []byte) {
	y, err := yaml.JSONToYAML(src)
	if err != nil {
		fmt.Printf("Error converting JSON into YAML: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(y))
}
