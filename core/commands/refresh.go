// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"dhcli/core"
	"dhcli/core/service"

	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh <environment>",
	Short: "Refresh access token",
	Long:  "Refresh the access token of a given environment.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var environment string
		if len(args) > 0 {
			environment = args[0]
		}

		service.RefreshHandler(environment)
	},
}

func init() {
	core.RegisterCommand(refreshCmd)
}
