package cmd

import (
	"log"

	"dhcli/service"
	"github.com/spf13/cobra"
)

var (
	envFlag     string
	outFlag     string
	projectFlag string
	nameFlag    string
)

var getCmd = &cobra.Command{
	Use:   "get <resource> [id]",
	Short: "Retrieve a DigitalHub resource",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := ""
		if len(args) > 1 {
			id = args[1]
		}

		err := service.GetResource(envFlag, outFlag, projectFlag, nameFlag, args[0], id, args[1:])
		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
	},
}

func init() {
	getCmd.Flags().StringVarP(&envFlag, "e", "e", "", "Environment")
	getCmd.Flags().StringVarP(&outFlag, "o", "o", "short", "Output format (short, json, yaml)")
	getCmd.Flags().StringVarP(&projectFlag, "p", "p", "", "Project")
	getCmd.Flags().StringVarP(&nameFlag, "n", "n", "", "Name")
	RegisterCommand(getCmd)
}
