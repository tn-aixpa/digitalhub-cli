package cmd

import (
	"log"

	"dhcli/service"
	"github.com/spf13/cobra"
)

var (
	listEnv     string
	listOutput  string
	listProject string
	listName    string
	listKind    string
	listState   string
)

var listCmd = &cobra.Command{
	Use:   "list <resource>",
	Short: "List DigitalHub resources",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.ListResources(
			listEnv, listOutput, listProject,
			listName, listKind, listState,
			args[0],
		); err != nil {
			log.Fatalf("List failed: %v", err)
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&listEnv, "e", "e", "", "environment")
	listCmd.Flags().StringVarP(&listOutput, "o", "o", "short", "output format: short|json|yaml")
	listCmd.Flags().StringVarP(&listProject, "p", "p", "", "project")
	listCmd.Flags().StringVarP(&listName, "n", "n", "", "name")
	listCmd.Flags().StringVarP(&listKind, "k", "k", "", "kind")
	listCmd.Flags().StringVarP(&listState, "s", "s", "", "state")

	RegisterCommand(listCmd)
}
