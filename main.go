package main

import (
	"dhcli/cmd"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	cmd.ExecuteCommand(os.Args[1:])
}
