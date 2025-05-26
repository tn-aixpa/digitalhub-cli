package cmd

import (
	"flag"
	"log"
	"os"

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
	if len(args) < 1 {
		log.Printf("Error: Environment name is a required positional argument.\nUsage: dhcli use <environment>\n")
		os.Exit(1)
	}

	environmentName := args[0]
	cfg := utils.LoadIni(false)
	if !cfg.HasSection(environmentName) {
		log.Printf("Specified environment does not exist.\n")
		os.Exit(1)
	}

	defaultSection := cfg.Section("DEFAULT")
	defaultSection.Key("current_environment").SetValue(environmentName)

	utils.SaveIni(cfg)
	log.Printf("Switched default to '%v'.\n", environmentName)
}
