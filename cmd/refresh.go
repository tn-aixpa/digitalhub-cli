package cmd

import (
	"dhcli/service"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh <environment>",
	Short: "Refresh environment variable",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.RefreshHandler(args[0])
	},
}

func init() {
	RegisterCommand(refreshCmd)
}
