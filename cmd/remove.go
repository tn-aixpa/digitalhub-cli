package cmd

import (
	"flag"
	"log"

	"gopkg.in/ini.v1"

	"dhcli/utils"
)

func init() {

	RegisterCommand(&Command{
		Name:        "remove",
		Description: "./dhcli remove <environment>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     removeHandler,
	})

}

func removeHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		log.Fatalf("Error: Environment is a required positional argument.\nUsage: ./dhcli remove <environment>")
	}

	sectionName := args[0]

	cfg := utils.LoadIni(false)

	cfg.DeleteSection(sectionName)

	defaultSection := cfg.Section("DEFAULT")
	if defaultSection.Key("current_environment").String() == sectionName {
		defaultSection.DeleteKey("current_environment")
	}

	utils.SaveIni(cfg)
	log.Printf("'%v' has been removed.", sectionName)
}
