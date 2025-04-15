package cmd

import (
	"dhcli/utils"
	"flag"
	"log"

	"gopkg.in/ini.v1"
)

func init() {
	RegisterCommand(&Command{
		Name:        "delete",
		Description: "dhcli delete [-n <name> -e <entity type> -i <id>] <project>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("n", "", "environment name")
			fs.String("e", "", "entity type")
			fs.String("i", "", "id")
		},
		Handler: deleteHandler,
	})
}

func deleteHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		log.Fatalf("Error: Project is required.")
	}
	fs.Parse(args)
	name := fs.Lookup("n").Value.String()
	project := fs.Args()[0]
	entityType := fs.Lookup("e").Value.String()
	id := fs.Lookup("i").Value.String()

	if entityType != "" && id == "" {
		log.Fatalf("Entity type specified, but ID missing.")
	} else if entityType == "" && id != "" {
		log.Fatalf("ID specified, but entity type missing.")
	}

	_, section := loadConfig([]string{name})

	method := "DELETE"
	url := utils.BuildCoreUrl(section, method, project, entityType, id)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())
	utils.DoRequest(req)
	log.Printf("Deleted successfully.")
}
