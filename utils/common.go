package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

func BuildCoreUrl(section *ini.Section, method string, project string, entity string, id string) string {
	base := section.Key("endpoint").String() + "/api/" + section.Key("api_version").String()
	mid := "/projects"
	endpoint := ""
	if project != "" {
		endpoint = "/" + project
		if entity != "" {
			mid = "/-"
			endpoint += "/" + entity
			if id != "" {
				if method != "DELETE" || entity == "secrets" || entity == "runs" {
					endpoint += "/" + id
				} else {
					endpoint += "?name=" + id
				}
			}
		} else if method == "DELETE" {
			endpoint += "?cascade=true"
		}
	}

	return base + mid + endpoint
}

func PrepareRequest(method string, url string, data []byte, accessToken string) *http.Request {
	var body io.Reader = nil
	if data != nil {
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Printf("Failed to initialize request: %v\n", err)
		os.Exit(1)
	}

	if data != nil {
		req.Header.Add("Content-type", "application/json")
	}

	if accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+accessToken)
	}

	return req
}

func DoRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error performing request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Core responded with error %v\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	return body, err
}

func DoRequestAndPrintResponse(req *http.Request) {
	body, err := DoRequest(req)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(body))
}
