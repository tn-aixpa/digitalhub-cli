package cmd

import (
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

		if err := service.RegisterEnvironment(envFlag, endpoint); err != nil {
			log.Fatalf("Registration failed: %v", err)
		}
	},
}

func init() {
	registerCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment")
	RegisterCommand(registerCmd)
}
