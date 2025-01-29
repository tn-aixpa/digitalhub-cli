package utils

import (
	"bufio"
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
	iniPath += string(os.PathSeparator) + iniName

	return iniPath
}

func LoadIni(createOnMissing bool) *ini.File {
	cfg, err := ini.Load(getIniPath())
	if err != nil {
		if !createOnMissing {
			log.Fatalf("Failed to read ini file: %v", err)
		}
		return ini.Empty()
	}

	return cfg
}

func SaveIni(cfg *ini.File) {
	err := cfg.SaveTo(getIniPath())
	if err != nil {
		log.Fatalf("Failed to update ini file: %v", err)
	}
}

func GitignoreAddIniFile() {
	path := "./.gitignore"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Cannot open .gitignore file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == iniName {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading .gitignore file contents: %v", err)
	}

	if _, err = f.WriteString(iniName); err != nil {
		log.Fatalf("Error while adding entry to .gitignore file: %v", err)
	}
}
