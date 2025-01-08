package cmd

import "flag"

func init() {

	RegisterCommand(&Command{
		Name:        "register",
		Description: "DH CLI register ",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     registerHandler,
	})

}

func registerHandler(args []string, fs *flag.FlagSet) {

}

/*

$: dhcli register bologna <url>

Scrivere file ini nella home personale altrimenti se non c'e' nella cartella corrente ( da dove lancio il comando)
~/mario
	.cli.ini

[bologna]
authorization_endpoint: <url>
client_id: xxxxxxx
token_endpoint: <url>

[roma]
....
[milano]


$: dhcli login bologna

quando ritorna il token aggiorno il file ini per la sezione bologna.

[bologna]
....
jwt_token: <token ricevuto dall'autorizzazione>
*/
