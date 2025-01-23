package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"gopkg.in/ini.v1"
)

func init() {
	RegisterCommand(&Command{
		Name:        "init",
		Description: "./dhcli init <endpoint_core>",
		SetupFlags:  func(fs *flag.FlagSet) {},
		Handler:     initHandler,
	})
}

func initHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		log.Fatalf("Error: Core endpoint is a required positional argument.\nUsage: ./dhcli init <core_endpoint>")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", args[0], nil)
	if err != nil {
		log.Fatalf("Failed to set up HTTP request: %v", err)
	}
	// TODO get token
	req.Header.Set("Authorization", "Bearer "+args[1])

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching API version: %v", err)
	}
	defer resp.Body.Close()

	apiVersion := resp.Header["X-Api-Version"][0]
	apiVersionMinor := apiVersion[:strings.LastIndex(apiVersion, ".")]

	cmd := "pip install digitalhub~=" + apiVersionMinor
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	fmt.Println(string(out))
}
