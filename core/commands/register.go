package commands

import (
	"dhcli/core"
	"log"

	"dhcli/service"
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
