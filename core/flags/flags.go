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

func AddCommonFlags(cmd *cobra.Command, flagsToAdd ...string) {

	if len(flagsToAdd) == 0 {
		flagsToAdd = []string{"env", "out", "project", "name"}
	}

	for _, flag := range flagsToAdd {
		switch flag {
		case "env":
			cmd.Flags().StringVarP(&EnvFlag, "env", "e", "", "environment")
		case "out":
			cmd.Flags().StringVarP(&OutFlag, "out", "o", "short", "output format (short, json, yaml)")
		case "project":
			cmd.Flags().StringVarP(&ProjectFlag, "project", "p", "", "project")
		case "name":
			cmd.Flags().StringVarP(&NameFlag, "name", "n", "", "name")
		}
	}
}
