package cmd

import (
	"dhcli/cmd/root"
	"dhcli/service"
	"github.com/spf13/cobra"
)

var listEnvCmd = &cobra.Command{
	Use:   "list-env",
	Short: "Fetch environment variables",
	Run: func(cmd *cobra.Command, args []string) {

		service.ListEnvHandler()

	},
}

func init() {
	root.RegisterCommand(listEnvCmd)
}
