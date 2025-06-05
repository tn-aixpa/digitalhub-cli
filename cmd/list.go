package cmd

import (
	"dhcli/utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strconv"

	"sigs.k8s.io/yaml"
)

func init() {
	RegisterCommand(&Command{
		Name:        "list",
		Description: "dhcli list [-e <environment> -o <output format> -p <project> -n <name> -k <kind> -s <state>] <resource>",
		SetupFlags: func(fs *flag.FlagSet) {
			// CLI-specific
			fs.String("e", "", "environment")
			fs.String("o", "", "output format")

			// API
			fs.String("p", "", "project")
			fs.String("n", "", "name")
			fs.String("k", "", "kind")
			fs.String("s", "", "state")
		},
		Handler: listHandler,
	})
}

func listHandler(args []string, fs *flag.FlagSet) {
	fs.Parse(args)
	if len(fs.Args()) < 1 {
		log.Println("Error: resource type is required.")
		os.Exit(1)
	}
	resource := utils.TranslateEndpoint(fs.Args()[0])

	// Load environment and check API level requirements
	environment := fs.Lookup("e").Value.String()
	cfg, section := loadIniConfig([]string{environment})
	utils.CheckUpdateEnvironment(cfg, section)
	utils.CheckApiLevel(section, utils.ListMin, utils.ListMax)

	format := utils.TranslateFormat(fs.Lookup("o").Value.String())
	project := fs.Lookup("p").Value.String()

	if resource != "projects" && project == "" {
		log.Println("Project is mandatory when performing this operation on resources other than projects.")
		os.Exit(1)
	}

	params := map[string]string{}
	params["name"] = fs.Lookup("n").Value.String()
	if params["name"] != "" {
		params["versions"] = "all"
	}
	params["kind"] = fs.Lookup("k").Value.String()
	params["state"] = fs.Lookup("s").Value.String()
	params["size"] = "200"
	params["sort"] = "updated,asc"

	method := "GET"
	url := utils.BuildCoreUrl(section, project, resource, "", params)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())

	body, err := utils.DoRequest(req)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	mapResp := map[string]interface{}{}
	json.Unmarshal([]byte(string(body)), &mapResp)

	elements := mapResp["content"].([]interface{})
	pageableMap := mapResp["pageable"].(map[string]interface{})
	pageNumber := int(reflect.ValueOf(pageableMap["pageNumber"]).Float())
	totalPages := int(reflect.ValueOf(mapResp["totalPages"]).Float())

	for {
		if pageNumber >= totalPages-1 {
			break
		}
		params["page"] = strconv.Itoa(pageNumber + 1)

		url := utils.BuildCoreUrl(section, project, resource, "", params)
		req = utils.PrepareRequest(method, url, nil, section.Key("access_token").String())

		body, err = utils.DoRequest(req)
		if err != nil {
			log.Printf("Error reading response: %v\n", err)
			os.Exit(1)
		}

		mapResp := map[string]interface{}{}
		json.Unmarshal([]byte(string(body)), &mapResp)

		elements = slices.Concat(elements, mapResp["content"].([]interface{}))
		pageableMap = mapResp["pageable"].(map[string]interface{})
		pageNumber = int(reflect.ValueOf(pageableMap["pageNumber"]).Float())
	}

	if format == "short" {
		printShortList(elements)
	} else if format == "json" {
		printJsonList(elements)
	} else if format == "yaml" {
		utils.PrintCommentForYaml(section, args)
		printYamlList(elements)
	}
}

func printShortList(resources []interface{}) {
	printShortLineList("NAME", "ID", "KIND", "UPDATED", "STATE", "LABELS")
	fmt.Println("")

	for _, r := range resources {
		m := r.(map[string]interface{})

		rName := m["name"].(string)
		rId := m["id"].(string)
		rKind := m["kind"].(string)
		rUpdated := ""
		rState := ""
		rLabels := ""

		metadata := m["metadata"]
		if reflect.ValueOf(metadata).Kind() == reflect.Map {
			metadataMap := metadata.(map[string]interface{})
			rUpdated = metadataMap["updated"].(string)
			labels := metadataMap["labels"]
			if labels != nil {
				labelsArray := labels.([]interface{})
				for _, label := range labelsArray {
					rLabels = rLabels + label.(string) + ", "
				}
				rLabels = rLabels[:len(rLabels)-2]
			}
		}

		status := m["status"]
		if reflect.ValueOf(status).Kind() == reflect.Map {
			statusMap := status.(map[string]interface{})
			rState = statusMap["state"].(string)
		}

		printShortLineList(rName, rId, rKind, rUpdated, rState, rLabels)
	}
}

func printShortLineList(rName string, rId string, rKind string, rUpdated string, rState string, rLabels string) {
	fmt.Printf("%-24s%-36s%-24s%-30s%-12s%s\n", rName, rId, rKind, rUpdated, rState, rLabels)
}

func printJsonList(src []interface{}) {
	j, err := json.MarshalIndent(src, "", "    ")
	if err != nil {
		log.Printf("Error while parsing resource array: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(j))
}

func printYamlList(src []interface{}) {
	y, err := yaml.Marshal(src)
	if err != nil {
		log.Printf("Error while parsing resource array: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(y))
}
