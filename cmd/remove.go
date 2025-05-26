package cmd

import (
	"flag"
	"log"
	"os"

	"dhcli/utils"
)

func init() {

	RegisterCommand(&Command{
		Name:        "remove",
		Description: "dhcli remove <environment>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     removeHandler,
	})

}

func removeHandler(args []string, fs *flag.FlagSet) {
	if len(args) < 1 {
		log.Printf("Error: Environment is a required positional argument.\nUsage: dhcli remove <environment>\n")
		os.Exit(1)
	}

	sectionName := args[0]

	cfg := utils.LoadIni(false)
	if !cfg.HasSection(sectionName) {
		log.Printf("Specified environment does not exist.\n")
		os.Exit(0)
	}

	cfg.DeleteSection(sectionName)

	defaultSection := cfg.Section("DEFAULT")
	if defaultSection.Key("current_environment").String() == sectionName {
		defaultSection.DeleteKey("current_environment")
	}

	utils.SaveIni(cfg)
	log.Printf("'%v' has been removed.\n", sectionName)
}
