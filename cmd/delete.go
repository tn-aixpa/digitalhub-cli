package cmd

import (
	"dhcli/utils"
	"flag"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

func init() {
	RegisterCommand(&Command{
		Name:        "delete",
		Description: "dhcli delete [-e <environment> -p <project> -c] <resource> <id>",
		SetupFlags: func(fs *flag.FlagSet) {
			// CLI-specific
			fs.String("e", "", "environment")

			// API
			fs.String("p", "", "project")
			fs.Bool("c", false, "cascade")
		},
		Handler: deleteHandler,
	})
}

func deleteHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	fs.Parse(args)
	if len(fs.Args()) < 2 {
		log.Println("Error: resource type and id are required.")
		os.Exit(1)
	}
	resource := utils.TranslateEndpoint(fs.Args()[0])
	id := fs.Args()[1]

	environment := fs.Lookup("e").Value.String()
	cascade := fs.Lookup("c").Value.String()
	project := fs.Lookup("p").Value.String()

	if resource != "projects" && project == "" {
		log.Println("Project is mandatory when performing this operation on resources other than projects.")
		os.Exit(1)
	}

	_, section := loadConfig([]string{environment})

	params := map[string]string{}
	params["cascade"] = cascade

	method := "DELETE"
	url := utils.BuildCoreUrl(section, project, resource, id, params)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())
	utils.DoRequest(req)
	log.Println("Deleted successfully.")
}
