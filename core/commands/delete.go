package commands

import (
	"dhcli/core"
	"dhcli/core/flags"
	"errors"
	"log"

	"dhcli/service"
	"github.com/spf13/cobra"
)

var (
	confirmFlag bool
	cascadeFlag bool
)
var deleteCmd = &cobra.Command{
	Use:   "delete <resource> [id]",
	Short: "Delete a resource from the core platform by ID or name",
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

		err := service.DeleteHandler(
			flags.EnvFlag,
			flags.ProjectFlag,
			flags.NameFlag,
			confirmFlag,
			cascadeFlag,
			args[0],
			id,
			args[1:])

		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(deleteCmd, "env", "project", "name")

	// Add file flags
	deleteCmd.Flags().BoolVarP(&confirmFlag, "confirm", "y", false, "skips the deletion confirmation prompt")
	deleteCmd.Flags().BoolVarP(&cascadeFlag, "cascade", "c", false, "if set, also deletes related resources (for projects)")

	core.RegisterCommand(deleteCmd)
}
