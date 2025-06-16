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

var getCmd = &cobra.Command{
	Use:   "get <resource> [id]",
	Short: "Retrieve a resource",
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

		err := service.GetHandler(
			flags.EnvFlag,
			flags.OutFlag,
			flags.ProjectFlag,
			flags.NameFlag,
			args[0],
			id,
			args[1:])

		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(getCmd)
	core.RegisterCommand(getCmd)
}
