// cmd/login.go
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"dhcli/service"
)

var loginCmd = &cobra.Command{
	Use:   "login [environment]",
	Short: "Login to a given environment",
	Long:  "Authenticates the user using OAuth2 PKCE flow with the specified environment.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var environment string
		if len(args) > 0 {
			environment = args[0]
		}

		if err := service.Login(environment); err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	},
}

func init() {
	RegisterCommand(loginCmd)
}
