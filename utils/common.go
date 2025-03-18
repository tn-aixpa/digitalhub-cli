package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
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
		log.Fatalf("Failed to initialize request: %v", err)
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
		log.Fatalf("Error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Core responded with error %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	return body, err
}

func DoRequestAndPrintResponse(req *http.Request) {
	body, err := DoRequest(req)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Println(string(body))
}
