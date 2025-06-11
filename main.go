package main

import (
	"dhcli/cmd"
	"gopkg.in/ini.v1"
)

func main() {
	ini.DefaultHeader = true
	cmd.Execute()
}
