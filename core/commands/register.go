// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"dhcli/core"
	"dhcli/core/service"
	"log"

	"github.com/spf13/cobra"
)

var (
	envFlag string
)

var registerCmd = &cobra.Command{
	Use:   "register <endpoint>",
	Short: "Register the configuration of a core instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := args[0]

		if err := service.RegisterHandler(envFlag, endpoint); err != nil {
			log.Fatalf("Registration failed: %v", err)
		}
	},
}

func init() {
	registerCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment")
	core.RegisterCommand(registerCmd)
}
