// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"dhcli/core"
	"dhcli/core/flags"
	"github.com/spf13/cobra"
)

var uploadFlag = flags.SpecificCommandFlag{}

var uploadCmd = &cobra.Command{
	Use:   "upload <resource>",
	Short: "upload a resource on the S3 aws",
	Long:  "Upload an artifact from ........................",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//if err := service.DownloadHandler(
		//	flags.EnvFlag,
		//	fileOrDirectoryFlag,
		//	flags.ProjectFlag,
		//	flags.NameFlag,
		//	args[0],
		//	id,
		//	args[1:]); err != nil {
		//	log.Fatalf("Download failed: %v", err)
		//}
	},
}

func init() {
	flags.AddCommonFlags(uploadCmd, "env", "project", "name")

	uploadCmd.Flags().StringVarP(&uploadFlag.InputFlag, "input", "i", "", "input filename or directory")
	uploadCmd.Flags().BoolVarP(&uploadFlag.CreateFlag, "create", "c", false, "if set, also create resource on core")

	core.RegisterCommand(uploadCmd)
}
