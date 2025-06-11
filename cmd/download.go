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
		id := ""
		if len(args) > 1 {
			id = args[1]
		}

		if err := service.DownloadFileWithOptions(
			flags.EnvFlag,
			flags.OutFlag,
			flags.ProjectFlag,
			flags.NameFlag,
			args[0],
			id,
			args[1:]); err != nil {
			_ = fmt.Errorf("download failed: %w", err)
		}
	},
}

func init() {
	flags.AddCommonFlags(downloadCmd)

	// override output common flag in this case out is a new filename or directory nam
	flag := downloadCmd.Flags().Lookup("out")

	if flag != nil {
		flag.Usage = "output filename or directory"
		flag.DefValue = "current filename or directory"
		err := flag.Value.Set("")
		if err != nil {
			return
		}
	}

	RegisterCommand(downloadCmd)
}
