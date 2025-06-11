package cmd

import (
	"log"

	"dhcli/service"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login [environment]",
	Short: "Authenticate using OAuth2 PKCE flow",
	Long:  "Initiates the OAuth2 PKCE login flow for the specified environment (or default).",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := ""
		if len(args) > 0 {
			env = args[0]
		}
		if err := service.Login(env); err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	},
}

func init() {
	RegisterCommand(loginCmd)
}
