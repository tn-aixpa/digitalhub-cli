package flags

import (
	"github.com/spf13/cobra"
)

var (
	EnvFlag     string
	OutFlag     string
	ProjectFlag string
	NameFlag    string
)

// AddCommonFlags aggiunge le flag comuni al comando dato
func AddCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&EnvFlag, "env", "e", "", "environment")
	cmd.Flags().StringVarP(&OutFlag, "out", "o", "short", "output format (short, json, yaml)")
	cmd.Flags().StringVarP(&ProjectFlag, "project", "p", "", "project")
	cmd.Flags().StringVarP(&NameFlag, "name", "n", "", "name")
}
