package main

import (
	"dhcli/cmd"
	"os"
)

func main() {
	cmd.ExecuteCommand(os.Args[1:]) //[]string{"login", "aac.digitalhub-dev.smartcommunitylab.it"}) //
}
