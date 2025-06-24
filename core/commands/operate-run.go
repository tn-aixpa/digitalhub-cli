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

var operateRunCmd = &cobra.Command{
	Use:   "operate-run <project> <id> <operation>",
	Short: "Perform an operation on a run",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		err := service.OperateRunHandler(
			flags.CommonFlag.EnvFlag,
			args[0],
			args[1],
			args[2])

		if err != nil {
			log.Fatalf("Failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(operateRunCmd, "env")

	core.RegisterCommand(operateRunCmd)
}
