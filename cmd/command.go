package cmd

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Description string
	SetupFlags  func(set *flag.FlagSet)
	Handler     func(args []string, fs *flag.FlagSet)
}

var commands = map[string]*Command{}

func RegisterCommand(cmd *Command) {
	commands[cmd.Name] = cmd
}

func ExecuteCommand(args []string) {
	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println("Usage: dhcli <command> [options]")
		fmt.Println("\nAvailable commands:")
		for _, cmd := range commands {
			fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
		}
		os.Exit(1)
	}

	// Extract the command
	commandName := args[0]
	command, exists := commands[commandName]
	if !exists {
		fmt.Printf("Unknown command: %s\n", commandName)
		fmt.Println("Run dhcli to see available commands.")
		os.Exit(1)
	}

	// Create a new FlagSet for this command
	fs := flag.NewFlagSet(command.Name, flag.ExitOnError)

	// Set up flags for the command
	command.SetupFlags(fs)

	// Execute the command handler
	command.Handler(args[1:], fs)
}
