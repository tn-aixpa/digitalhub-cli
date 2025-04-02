package utils

import (
	"fmt"
	"os"
	"reflect"
	"time"

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
			fmt.Printf("Failed to read ini file: %v\n", err)
			os.Exit(1)
		}
		return ini.Empty()
	}

	return cfg
}

func SaveIni(cfg *ini.File) {
	err := cfg.SaveTo(getIniPath())
	if err != nil {
		fmt.Printf("Failed to update ini file: %v\n", err)
		os.Exit(1)
	}
}

func ReflectValue(v interface{}) string {
	f := reflect.ValueOf(v)

	switch f.Kind() {
	case reflect.String:
		return f.String()
	case reflect.Int, reflect.Int64:
		return fmt.Sprint(f.Int())
	case reflect.Uint, reflect.Uint64:
		return fmt.Sprint(f.Uint())
	case reflect.Float64:
		return fmt.Sprint(f.Float())
	case reflect.Bool:
		return fmt.Sprint(f.Bool())
	case reflect.TypeOf(time.Now()).Kind():
		return f.Interface().(time.Time).Format(time.RFC3339)
	case reflect.Slice:
		return fmt.Sprint(f.Interface())
	default:
		return ""
	}
}
