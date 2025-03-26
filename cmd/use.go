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
		Name:        "use",
		Description: "dhcli use <environment>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     useHandler,
	})

}

func useHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		fmt.Printf("Error: Environment name is a required positional argument.\nUsage: dhcli use <environment>")
		os.Exit(1)
	}

	environmentName := args[0]
	cfg := utils.LoadIni(false)
	if !cfg.HasSection(environmentName) {
		fmt.Printf("Specified environment does not exist.")
		os.Exit(1)
	}

	defaultSection := cfg.Section("DEFAULT")
	defaultSection.Key("current_environment").SetValue(environmentName)

	utils.SaveIni(cfg)
	fmt.Printf("Switched default to '%v'.", environmentName)
}
