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
		Name:        "remove",
		Description: "dhcli remove <environment>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     removeHandler,
	})

}

func removeHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		fmt.Printf("Error: Environment is a required positional argument.\nUsage: dhcli remove <environment>")
		os.Exit(1)
	}

	sectionName := args[0]

	cfg := utils.LoadIni(false)
	if !cfg.HasSection(sectionName) {
		fmt.Printf("Specified environment does not exist.")
		os.Exit(0)
	}

	cfg.DeleteSection(sectionName)

	defaultSection := cfg.Section("DEFAULT")
	if defaultSection.Key("current_environment").String() == sectionName {
		defaultSection.DeleteKey("current_environment")
	}

	utils.SaveIni(cfg)
	fmt.Printf("'%v' has been removed.", sectionName)
}
