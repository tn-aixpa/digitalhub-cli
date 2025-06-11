package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"dhcli/utils"
)

func InitEnvironment(env string, pre bool) error {
	// Check Python version
	versionOutput, err := exec.Command("python3", "--version").Output()
	if err != nil {
		return fmt.Errorf("python3 not installed or not in PATH: %w", err)
	}

	if !supportedPythonVersion(string(versionOutput)) {
		return fmt.Errorf("unsupported Python version (must be 3.9 <= v <= 3.12): %s", string(versionOutput))
	}

	_, section := utils.LoadIniConfig([]string{env})

	apiVersionMinor := section.Key("dhcore_version").String()
	versionParts := strings.SplitN(apiVersionMinor, ".", 3)
	if len(versionParts) >= 2 {
		apiVersionMinor = strings.Join(versionParts[:2], ".")
	}

	// Prompt user
	if !confirm(fmt.Sprintf("Newest patch version of DigitalHub %s will be installed, continue? Y/n", apiVersionMinor)) {
		log.Println("Installation cancelled by user.")
		return nil
	}

	pipSpecifier := "~=" + apiVersionMinor + ".0"
	for _, pkg := range packageList() {
		args := []string{"-m", "pip", "install", pkg + pipSpecifier}
		if pre {
			args = append(args, "--pre")
		}

		cmd := exec.Command("python3", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		log.Printf("Installing %s...\n", pkg+pipSpecifier)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("pip install failed for %s: %w", pkg, err)
		}
	}

	log.Println("Installation complete.")
	return nil
}

func supportedPythonVersion(version string) bool {
	version = strings.TrimSpace(version)
	if !strings.HasPrefix(version, "Python ") {
		return false
	}

	ver := strings.TrimPrefix(version, "Python ")
	parts := strings.Split(ver, ".")
	if len(parts) < 2 {
		return false
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil || major != 3 {
		return false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil || minor < 9 || minor > 12 {
		return false
	}

	return true
}

func confirm(prompt string) bool {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			os.Exit(1)
		}
		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y", "":
			return true
		case "n":
			return false
		default:
			fmt.Print("Please answer y or n: ")
		}
	}
}

func packageList() []string {
	return []string{
		"digitalhub[full]",
		"digitalhub-runtime-python",
	}
}
