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
		Name:        "remove",
		Description: "DH CLI remove",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     removeHandler,
	})

}

func removeHandler(args []string, fs *flag.FlagSet) {
	if len(args) < 1 {
		log.Fatalf("Error: Environment is a required positional argument.\nUsage: dhcli remove <environment>")
	}

	sectionName := args[0]
	iniPath := getIniPath()

	cfg, err := ini.Load(iniPath)
	if err != nil {
		fmt.Printf("Failed to load ini file: %v", err)
		os.Exit(1)
	}

	cfg.DeleteSection(sectionName)

	err = cfg.SaveTo(iniPath)
	if err != nil {
		fmt.Printf("Failed to write ini file: %v", err)
		os.Exit(1)
	}
}
