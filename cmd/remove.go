package cmd

import (
	"dhcli/cmd/root"
	"dhcli/service"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <environment>",
	Short: "Remove an environment from the config",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.RemoveHandler(args[0])
	},
}

func init() {
	root.RegisterCommand(removeCmd)
}
