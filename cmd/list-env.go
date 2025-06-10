// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"dhcli/utils"
	"flag"
	"fmt"
	"log"
)

func init() {

	RegisterCommand(&Command{
		Name:        "list-env",
		Description: "dhcli list-env",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     listenvHandler,
	})

}

func listenvHandler(args []string, fs *flag.FlagSet) {
	cfg := utils.LoadIni(true)
	sections := cfg.SectionStrings()
	sectionsString := ""

	for _, name := range sections {
		if name != "DEFAULT" {
			sectionsString += name + ", "
		}
	}

	if sectionsString == "" {
		log.Println("No environments available.")
		return
	}
	sectionsString = sectionsString[:len(sectionsString)-2]

	log.Println("Available environments:")
	fmt.Printf("%v\n", sectionsString)
}
