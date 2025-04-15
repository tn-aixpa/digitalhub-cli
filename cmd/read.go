package cmd

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/ini.v1"

	"dhcli/utils"
)

func init() {
	RegisterCommand(&Command{
		Name:        "read",
		Description: "dhcli read [-n <name> -p <project> -e <entity type> -i <id>]",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("n", "", "environment name")
			fs.String("p", "", "project")
			fs.String("e", "", "entity type")
			fs.String("i", "", "id")
		},
		Handler: readHandler,
	})
}

func readHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	fs.Parse(args)
	name := fs.Lookup("n").Value.String()
	project := fs.Lookup("p").Value.String()
	entityType := fs.Lookup("e").Value.String()
	id := fs.Lookup("i").Value.String()

	if id != "" && (project == "" || entityType == "") {
		fmt.Println("ID specified, but project or entity type are missing.")
		os.Exit(1)
	}
	if project == "" && entityType != "" {
		fmt.Println("Entity type specified, but project is missing.")
		os.Exit(1)
	}

	_, section := loadConfig([]string{name})

	method := "GET"
	url := utils.BuildCoreUrl(section, method, project, entityType, id)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())
	utils.DoRequestAndPrintResponse(req)
}
