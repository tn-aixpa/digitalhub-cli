package cmd

import (
	"dhcli/cmd/flags"
	"log"

	"dhcli/service"
	"github.com/spf13/cobra"
)

var (
	listKind  string
	listState string
)

var listCmd = &cobra.Command{
	Use:   "list <resource>",
	Short: "List digitalHub resources",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.ListResources(
			flags.EnvFlag,
			flags.OutFlag,
			flags.ProjectFlag,
			flags.NameFlag,
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

	RegisterCommand(listCmd)
}
