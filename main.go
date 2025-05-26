package main

import (
	"dhcli/cmd"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

func main() {
	log.SetFlags(0)
	ini.DefaultHeader = true
	cmd.ExecuteCommand(os.Args[1:])
}
