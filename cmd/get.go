package cmd

import (
	"dhcli/cmd/flags"
	"log"

	"dhcli/service"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <resource> [id]",
	Short: "Retrieve a DigitalHub resource",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := ""
		if len(args) > 1 {
			id = args[1]
		}

		err := service.GetResource(
			flags.EnvFlag,
			flags.OutFlag,
			flags.ProjectFlag,
			flags.NameFlag,
			args[0],
			id,
			args[1:])

		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(getCmd)
	RegisterCommand(getCmd)
}
