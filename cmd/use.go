package cmd

import (
	"flag"
	"log"

	"gopkg.in/ini.v1"

	"dhcli/utils"
)

func init() {

	RegisterCommand(&Command{
		Name:        "use",
		Description: "./dhcli use <environment>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     useHandler,
	})

}

func useHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		log.Fatalf("Error: Environment name is a required positional argument.\nUsage: ./dhcli use <environment>")
	}

	environmentName := args[0]
	cfg := utils.LoadIni(false)
	if !cfg.HasSection(environmentName) {
		log.Fatalf("Specified environment does not exist.")
	}

	defaultSection := cfg.Section("DEFAULT")
	defaultSection.Key("current_environment").SetValue(environmentName)

	utils.SaveIni(cfg)
}
