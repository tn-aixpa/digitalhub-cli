package utils

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

const iniName = ".cli.ini"

func getIniPath() string {
	iniPath, err := os.UserHomeDir()
	if err != nil {
		iniPath = "."
	}
	iniPath += "/" + iniName

	return iniPath
}

func LoadIni() *ini.File {
	cfg, err := ini.Load(getIniPath())
	if err != nil {
		log.Fatalf("Failed to read ini file: %v", err)
	}

	return cfg
}

func SaveIni(cfg *ini.File) {
	err := cfg.SaveTo(getIniPath())
	if err != nil {
		log.Fatalf("Failed to update ini file: %v", err)
	}
}
