// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"dhcli/core"
	"dhcli/core/flags"
	"dhcli/core/service"
	"log"

	"github.com/spf13/cobra"
)

var (
	filePathFlag string
	resetIdFlag  bool
)
var createCmd = &cobra.Command{
	Use:   "create <resource>",
	Short: "Creates a new resource from a YAML file (or a name for projects)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := service.CreateHandler(
			flags.EnvFlag,
			flags.ProjectFlag,
			flags.NameFlag,
			filePathFlag,
			resetIdFlag,
			args[0])
		if err != nil {
			log.Fatalf("Create failed: %v", err)
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
