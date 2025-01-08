package cmd

import "flag"

func init() {

	RegisterCommand(&Command{
		Name:        "logout",
		Description: "DH CLI logout",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     logoutHandler,
	})

}

func logoutHandler(args []string, fs *flag.FlagSet) {

}
