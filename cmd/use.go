package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

func init() {

	RegisterCommand(&Command{
		Name:        "use",
		Description: "DH CLI use",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     useHandler,
	})

}

func useHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		log.Fatalf("Error: Environment name is a required positional argument.\nUsage: dhcli use <environment>")
	}

	environmentName := args[0]
	iniPath := getIniPath()

	cfg, err := ini.Load(iniPath)
	if err != nil {
		fmt.Printf("Failed to load ini file: %v", err)
		os.Exit(1)
	}

	if !cfg.HasSection(environmentName) {
		fmt.Printf("Specified environment does not exist.")
		os.Exit(1)
	}

	defaultSection := cfg.Section("DEFAULT")
	defaultSection.Key("environment").SetValue(environmentName)

	err = cfg.SaveTo(iniPath)
	if err != nil {
		fmt.Printf("Failed to write ini file: %v", err)
		os.Exit(1)
	}
}
