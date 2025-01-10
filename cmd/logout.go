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
		Name:        "logout",
		Description: "DH CLI logout ",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     logoutHandler,
	})

}

func logoutHandler(args []string, fs *flag.FlagSet) {
	if len(args) < 1 {
		log.Fatalf("Error: Name of configuration is a required positional argument.\nUsage: dhcli logout <name>")
	}

	sectionName := args[0]
	iniPath := getIniPath()

	cfg, err := ini.Load(iniPath)
	if err != nil {
		fmt.Printf("Failed to load ini file: %v", err)
		os.Exit(1)
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		log.Fatalf("Failed to read section '%s': %v.", sectionName, err)
	}

	section.DeleteKey("jwt_token")

	err = cfg.SaveTo(iniPath)
	if err != nil {
		fmt.Printf("Failed to write ini file: %v", err)
		os.Exit(1)
	}
}
