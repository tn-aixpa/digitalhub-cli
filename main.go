package main

import (
	"dhcli/core"
	_ "dhcli/core/commands"
	"gopkg.in/ini.v1"
)

func main() {
	ini.DefaultHeader = true
	core.Execute()
}
