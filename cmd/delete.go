package cmd

import (
	"dhcli/utils"
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	RegisterCommand(&Command{
		Name:        "delete",
		Description: "dhcli delete [-e <environment> -p <project> -n <name> -c -y] <resource> <id>",
		SetupFlags: func(fs *flag.FlagSet) {
			// CLI-specific
			fs.String("e", "", "environment")
			fs.Bool("y", false, "skip confirmation")

			// API
			fs.String("p", "", "project")
			fs.String("n", "", "name")
			fs.Bool("c", false, "cascade")
		},
		Handler: deleteHandler,
	})
}

func deleteHandler(args []string, fs *flag.FlagSet) {
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

	// Load environment and check API level requirements
	environment := fs.Lookup("e").Value.String()
	_, section := loadIniConfig([]string{environment})
	utils.CheckApiLevel(section, utils.DeleteMin, utils.DeleteMax)

	project := fs.Lookup("p").Value.String()
	name := fs.Lookup("n").Value.String()
	skipConfirm := fs.Lookup("y").Value.String()

	cascade := fs.Lookup("c").Value.String()
	if resource == "projects" && cascade != "true" {
		log.Println("WARNING: You are deleting a project without the cascade (-c) flag. Resources belonging to the project will not be deleted.")
	}

	if resource != "projects" && project == "" {
		log.Println("Project is mandatory when performing this operation on resources other than projects.")
		os.Exit(1)
	}

	params := map[string]string{}
	params["cascade"] = cascade

	confirmationMessage := fmt.Sprintf("Resource %v (%v) will be deleted, proceed? Y/n", id, resource)
	if id == "" {
		if name == "" {
			log.Println("You must specify id or name.")
			os.Exit(1)
		}
		if resource != "projects" {
			confirmationMessage = fmt.Sprintf("All versions of resource named '%v' (%v) will be deleted, proceed? Y/n", name, resource)
			params["name"] = name
		} else {
			confirmationMessage = fmt.Sprintf("Resource %v (%v) will be deleted, proceed? Y/n", name, resource)
			id = name
		}
	}

	// Ask for confirmation
	if skipConfirm != "true" {
		utils.WaitForConfirmation(confirmationMessage)
	}

	method := "DELETE"
	url := utils.BuildCoreUrl(section, project, resource, id, params)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())
	utils.DoRequest(req)
	log.Println("Deleted successfully.")
}
