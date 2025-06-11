package cmd

import (
	"dhcli/cmd/flags"
	"dhcli/service"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <resource> <id>",
	Short: "Download a resource from the S3 aws",
	Long:  "Download an artifact from ........................",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return errors.New("requires 1 or 2 arguments: <resource> [<id>]")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		id := ""
		if len(args) > 1 {
			id = args[1]
		}

		opts := service.DownloadOptions{
			Environment: flags.EnvFlag,
			Output:      flags.OutFlag,
			Project:     flags.ProjectFlag,
			Name:        flags.NameFlag,
			Resource:    resource,
			ID:          id,
		}

		if err := service.DownloadFileWithOptions(opts); err != nil {
			_ = fmt.Errorf("download failed: %w", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(downloadCmd)
	RegisterCommand(downloadCmd)
}
