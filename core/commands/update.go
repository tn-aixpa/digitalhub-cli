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
	fileFlag string
)

var updateCmd = &cobra.Command{
	Use:   "update <resource> <id>",
	Short: "Update a specific resource using data from a YAML file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := service.UpdateHandler(
			flags.CommonFlag.EnvFlag,
			flags.CommonFlag.ProjectFlag,
			fileFlag,
			args[0],
			args[1])

		if err != nil {
			log.Fatalf("Update failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(updateCmd, "env", "project")

	// Add file flags
	updateCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "path to the YAML file containing the resource data to be updated")
	core.RegisterCommand(updateCmd)
}
