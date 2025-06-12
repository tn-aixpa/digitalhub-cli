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
	Use:   "register [-e <environment>] <endpoint>",
	Short: "Register a DigitalHub core endpoint and store its configuration",
	Args:  cobra.ExactArgs(1), // endpoint richiesto
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
