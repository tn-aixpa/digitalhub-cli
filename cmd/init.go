package cmd

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	RegisterCommand(&Command{
		Name:        "init",
		Description: "dhcli init <environment>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.Bool("pre", false, "pip --pre flag")
		},
		Handler: initHandler,
	})
}

func initHandler(args []string, fs *flag.FlagSet) {
	fs.Parse(args)

	// Check if Python version is supported
	versionOutput, err := exec.Command("python3", "--version").Output()
	if err != nil {
		log.Printf("python3 does not seem to be installed: %v\n", err)
		os.Exit(1)
	}
	if !supportedPythonVersion(string(versionOutput)) {
		log.Printf("Python version is not supported (must be 3.9.xx <= v <=3.12.xx): %v\n", string(versionOutput))
		os.Exit(1)
	}

	// Read config from ini file
	loadArgs := args
	if len(args) > 0 && args[0] == "--pre" {
		loadArgs = args[1:]
	}
	_, section := loadIniConfig(loadArgs)

	apiVersionMinor := section.Key("dhcore_version").String()
	versionSplits := strings.SplitN(apiVersionMinor, ".", 3)
	if len(versionSplits) > 2 {
		apiVersionMinor = strings.Join(versionSplits[:2], ".")
	}

	// Ask for confirmation
	for {
		buf := bufio.NewReader(os.Stdin)
		log.Printf("Newest patch version of digitalhub %v will be installed, continue? Y/n\n", apiVersionMinor)
		userInput, err := buf.ReadBytes('\n')
		if err != nil {
			log.Printf("Error in reading user input: %v\n", err)
			os.Exit(1)
		} else {
			yn := strings.TrimSpace(string(userInput))
			if strings.ToLower(yn) == "y" || yn == "" {
				break
			} else if strings.ToLower(yn) == "n" {
				log.Println("Cancelling installation.")
				return
			}
			log.Println("Invalid input, must be y or n")
		}
	}

	pipOption := "~=" + apiVersionMinor + ".0"
	pre := fs.Lookup("pre").Value.String()

	for _, pkg := range packageList() {
		var cmd *exec.Cmd
		if pre != "true" {
			cmd = exec.Command("python3", "-m", "pip", "install", pkg+pipOption)
		} else {
			cmd = exec.Command("python3", "-m", "pip", "install", pkg+pipOption, "--pre")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to execute command: %v; %v\n", err)
			os.Exit(1)
		}
		log.Println("Installation complete.")
	}
}

func supportedPythonVersion(pythonVersion string) bool {
	// Remove 'Python ' part
	ver := pythonVersion[strings.Index(pythonVersion, " ")+1:]

	// Split and check major and minor
	majorString := ver[:strings.Index(ver, ".")]
	if majorString != "3" {
		return false
	}

	minorString := ver[strings.Index(ver, ".")+1 : strings.LastIndex(ver, ".")]
	minor, err := strconv.Atoi(minorString)
	if err != nil || minor < 9 || minor > 12 {
		return false
	}

	return true
}

func packageList() []string {
	return []string{
		"digitalhub[full]",
		"digitalhub-runtime-python",
	}
}
