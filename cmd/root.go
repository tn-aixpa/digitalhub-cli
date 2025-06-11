package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dhcli",
	Short: "dhcli is a tool for managing resource in core platform",
	Long:  `dhcli is a command-line utility for downloading, uploading, and managing core platform entity`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func RegisterCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
