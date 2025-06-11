package cmd

import (
	"log"

	"dhcli/service"
	"github.com/spf13/cobra"
)

//type OpenIDConfig struct {
//	AuthorizationEndpoint string   `json:"authorization_endpoint" ini:"authorization_endpoint"`
//	TokenEndpoint         string   `json:"token_endpoint" ini:"token_endpoint"`
//	Issuer                string   `json:"issuer" ini:"issuer"`
//	ClientID              string   `json:"dhcore_client_id" ini:"client_id"`
//	Scope                 []string `json:"scopes_supported" ini:"scopes_supported"`
//	AccessToken           string   `json:"access_token" ini:"access_token"`
//	RefreshToken          string   `json:"refresh_token" ini:"refresh_token"`
//}
//
//type CoreConfig struct {
//	Name     string `json:"dhcore_name" ini:"dhcore_name"`
//	Issuer   string `json:"issuer" ini:"issuer"`
//	Version  string `json:"dhcore_version" ini:"dhcore_version"`
//	ClientID string `json:"dhcore_client_id" ini:"client_id"`
//}

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
