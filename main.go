package main

import (
	"dhcli/cmd/root"
	"gopkg.in/ini.v1"
)

func main() {
	ini.DefaultHeader = true
	root.Execute()
}
