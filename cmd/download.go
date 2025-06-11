package cmd

import (
	"dhcli/service"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func DownloadCmd() *cobra.Command {
	var env, output, project, name string

	cmd := &cobra.Command{
		Use:   "download <resource> <id>",
		Short: "Download a resource from the cloud",
		Long:  "dhcli download [-e <environment> -o <output> -p <project> -n <name>] <resource> <id>]",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return errors.New("requires 1 or 2 arguments: <resource> [<id>]")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			resource := args[0]
			id := ""
			if len(args) > 1 {
				id = args[1]
			}

			opts := service.DownloadOptions{
				Environment: env,
				Output:      output,
				Project:     project,
				Name:        name,
				Resource:    resource,
				ID:          id,
			}

			if err := service.DownloadFileWithOptions(opts); err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&env, "environment", "e", "", "environment")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	cmd.Flags().StringVarP(&project, "project", "p", "", "project")
	cmd.Flags().StringVarP(&name, "name", "n", "", "name")

	return cmd
}

func init() {
	RegisterCommand(DownloadCmd())
}
