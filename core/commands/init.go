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

var initFlag = flags.SpecificCommandFlag{}

var initCmd = &cobra.Command{
	Use:   "init [<environment>]",
	Short: "Install python packages for an environment",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := ""
		if len(args) > 0 {
			env = args[0]
		}
		if err := service.InitEnvironmentHandler(env, initFlag.PreFlag); err != nil {
			log.Fatalf("Init failed: %v", err)
		}
	},
}

func init() {
	initCmd.Flags().BoolVar(&initFlag.PreFlag, "pre", false, "Include pre-release versions when installing")
	core.RegisterCommand(initCmd)
}
