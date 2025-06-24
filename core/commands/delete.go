// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"dhcli/core"
	"dhcli/core/flags"
	"dhcli/core/service"
	"errors"
	"log"

	"github.com/spf13/cobra"
)

var deleteFlag = flags.SpecificCommandFlag{}

var deleteCmd = &cobra.Command{
	Use:   "delete <resource> [id]",
	Short: "Delete a resource by ID or name",
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
			flags.CommonFlag.EnvFlag,
			flags.CommonFlag.ProjectFlag,
			flags.CommonFlag.NameFlag,
			deleteFlag.ConfirmFlag,
			deleteFlag.CascadeFlag,
			args[0],
			id)

		if err != nil {
			log.Fatalf("Delete failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(deleteCmd, "env", "project", "name")

	// Add file flags
	deleteCmd.Flags().BoolVarP(&deleteFlag.ConfirmFlag, "confirm", "y", false, "skips the deletion confirmation prompt")
	deleteCmd.Flags().BoolVarP(&deleteFlag.CascadeFlag, "cascade", "c", false, "if set, also deletes related resources (for projects)")

	core.RegisterCommand(deleteCmd)
}
