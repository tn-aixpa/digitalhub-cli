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
	listKind  string
	listState string
)

var listCmd = &cobra.Command{
	Use:   "list <resource>",
	Short: "List resources",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.ListResourcesHandler(
			flags.CommonFlag.EnvFlag,
			flags.CommonFlag.OutFlag,
			flags.CommonFlag.ProjectFlag,
			flags.CommonFlag.NameFlag,
			listKind,
			listState,
			args[0],
		); err != nil {
			log.Fatalf("List failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(listCmd)

	// Add specific command flags
	listCmd.Flags().StringVarP(&listKind, "kind", "k", "", "kind")
	listCmd.Flags().StringVarP(&listState, "state", "s", "", "state")

	core.RegisterCommand(listCmd)
}
