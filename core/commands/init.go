package commands

import (
	"dhcli/core"
	"dhcli/core/service"
	"log"

	"github.com/spf13/cobra"
)

var preFlag bool

var initCmd = &cobra.Command{
	Use:   "init [<environment>]",
	Short: "Install digitalHub python packages for an environment",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := ""
		if len(args) > 0 {
			env = args[0]
		}
		if err := service.InitEnvironmentHandler(env, preFlag); err != nil {
			log.Fatalf("Init failed: %v", err)
		}
	},
}

func init() {
	initCmd.Flags().BoolVar(&preFlag, "pre", false, "Include pre-release versions when installing")
	core.RegisterCommand(initCmd)
}
