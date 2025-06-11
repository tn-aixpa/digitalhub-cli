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
	filePathFlag string
	resetIdFlag  bool
)
var createCmd = &cobra.Command{
	Use:   "create <resource> [id]",
	Short: "Creates a new resource on the core platform using data from a YAML file or a name",
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

		err := service.CreateHandler(
			flags.EnvFlag,
			flags.ProjectFlag,
			flags.NameFlag,
			filePathFlag,
			resetIdFlag,
			args[0],
			id,
			args[1:])

		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(createCmd, "env", "project", "name")

	// Add file flags
	createCmd.Flags().BoolVarP(&resetIdFlag, "reset-id", "r", false, "if set, removes the id field from the file to ensure the server assigns a new one")
	createCmd.Flags().StringVarP(&filePathFlag, "file", "f", "", "path to a YAML file containing the resource definition")

	core.RegisterCommand(createCmd)
}
