package cmd

import (
	"dhcli/service"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <environment>",
	Short: "Set and save the specified environment as the current one",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.UseHandler(args[0])
	},
}

func init() {
	RegisterCommand(useCmd)
}
