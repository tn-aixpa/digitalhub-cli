package utils

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

const (
	IniName            = ".dhcore.ini"
	CurrentEnvironment = "current_environment"
)

func getIniPath() string {
	iniPath, err := os.UserHomeDir()
	if err != nil {
		iniPath = "."
	}
	iniPath += string(os.PathSeparator) + IniName

	return iniPath
}

func LoadIni(createOnMissing bool) *ini.File {
	cfg, err := ini.Load(getIniPath())
	if err != nil {
		if !createOnMissing {
			fmt.Printf("Failed to read ini file: %v", err)
			os.Exit(1)
		}
		return ini.Empty()
	}

	return cfg
}

func SaveIni(cfg *ini.File) {
	err := cfg.SaveTo(getIniPath())
	if err != nil {
		fmt.Printf("Failed to update ini file: %v", err)
		os.Exit(1)
	}
}
