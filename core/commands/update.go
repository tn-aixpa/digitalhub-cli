package commands

import (
	"dhcli/core"
	"dhcli/core/flags"
	"dhcli/core/service"
	"errors"
	"log"

	"github.com/spf13/cobra"
)

var (
	fileFlag string
)

var updateCmd = &cobra.Command{
	Use:   "update <resource> [id]",
	Short: "Updates a specific resource using data from a YAML file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return errors.New("requires 1 or 2 arguments: <resource> [<id>]")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := ""
		if len(args) > 1 {
			id = args[1]
		}

		err := service.UpdateHandler(
			flags.EnvFlag,
			flags.ProjectFlag,
			fileFlag,
			args[0],
			id,
			args[1:])

		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(updateCmd, "env", "project")

	// Add file flags
	updateCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "path to the YAML file containing the resource data to be updated")
	core.RegisterCommand(updateCmd)
}
