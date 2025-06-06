package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type OpenIDConfig struct {
	AuthorizationEndpoint string   `json:"authorization_endpoint" ini:"authorization_endpoint"`
	TokenEndpoint         string   `json:"token_endpoint" ini:"token_endpoint"`
	Issuer                string   `json:"issuer" ini:"issuer"`
	ClientID              string   `json:"dhcore_client_id" ini:"client_id"`
	Scope                 []string `json:"scopes_supported" ini:"scopes_supported"`
	AccessToken           string   `json:"access_token" ini:"access_token"`
	RefreshToken          string   `json:"refresh_token" ini:"refresh_token"`
}

type CoreConfig struct {
	Name     string `json:"dhcore_name" ini:"dhcore_name"`
	Issuer   string `json:"issuer" ini:"issuer"`
	Version  string `json:"dhcore_version" ini:"dhcore_version"`
	ClientID string `json:"dhcore_client_id" ini:"client_id"`
}

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
			log.Printf("Failed to read ini file: %v\n", err)
			os.Exit(1)
		}
		return ini.Empty()
	}

	return cfg
}

func SaveIni(cfg *ini.File) {
	err := cfg.SaveTo(getIniPath())
	if err != nil {
		log.Printf("Failed to update ini file: %v\n", err)
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
		s := []string{}
		for _, element := range f.Interface().([]interface{}) {
			if reflect.ValueOf(element).Kind() == reflect.String {
				s = append(s, element.(string))
			}
		}
		return strings.Join(s, ",")
	default:
		return ""
	}
}

func BuildCoreUrl(section *ini.Section, project string, resource string, id string, params map[string]string) string {
	base := section.Key(DhCoreEndpoint).String() + "/api/" + section.Key("dhcore_api_version").String()
	endpoint := ""
	paramsString := ""
	if resource != "projects" && project != "" {
		endpoint += "/-/" + project
	}
	endpoint += "/" + resource
	if id != "" {
		endpoint += "/" + id
	}
	if params != nil && len(params) > 0 {
		paramsString = "?"
		for key, val := range params {
			if val != "" {
				paramsString += key + "=" + val + "&"
			}
		}
		paramsString = paramsString[:len(paramsString)-1]
	}

	return base + endpoint + paramsString
}

func PrepareRequest(method string, url string, data []byte, accessToken string) *http.Request {
	var body io.Reader = nil
	if data != nil {
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("Failed to initialize request: %v\n", err)
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
		log.Printf("Error performing request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Core responded with error %v\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	return body, err
}

func TranslateFormat(format string) string {
	lower := strings.ToLower(format)
	if lower == "json" {
		return "json"
	} else if lower == "yaml" || lower == "yml" {
		return "yaml"
	}
	return "short"
}

func loadConfig() map[string]interface{} {
	file, err := os.ReadFile("./" + configFile)
	if err != nil {
		log.Printf("Failed to read config file, some functionalities may not work: %v\n", err)
		return nil
	}

	var config map[string]interface{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Printf("Error unmarshalling config file, some functionalities may not work: %v\n", err)
		return nil
	}

	return config
}

func TranslateEndpoint(resource string) string {
	config := loadConfig()

	if config != nil {
		if endpoints, ok := config["resources"]; ok && reflect.ValueOf(endpoints).Kind() == reflect.Map {
			endpointsMap := endpoints.(map[string]interface{})

			for key, val := range endpointsMap {
				if key == resource {
					return val.(string)
				}

				if reflect.ValueOf(val).Kind() == reflect.String && resource == val.(string) {
					return resource
				}
			}
		}
	}

	log.Printf("Resource '%v' is not supported or the configuration file is invalid. Check or edit supported resources in %v.\n", resource, configFile)
	os.Exit(1)
	return ""
}

func WaitForConfirmation(msg string) {
	for {
		buf := bufio.NewReader(os.Stdin)
		log.Printf(msg)
		userInput, err := buf.ReadBytes('\n')
		if err != nil {
			log.Printf("Error in reading user input: %v\n", err)
			os.Exit(1)
		} else {
			yn := strings.TrimSpace(string(userInput))
			if strings.ToLower(yn) == "y" || yn == "" {
				break
			} else if strings.ToLower(yn) == "n" {
				log.Println("Cancelling.")
				os.Exit(0)
			}
			log.Println("Invalid input, must be y or n")
		}
	}
}

func PrintCommentForYaml(section *ini.Section, args []string) {
	fmt.Printf("# Generated on: %v\n", time.Now().Round(0))
	fmt.Printf("#   from environment: %v (core version %v)\n", section.Key("dhcore_name").String(), section.Key("dhcore_version").String())
	fmt.Printf("#   found at: %v\n", section.Key(DhCoreEndpoint).String())
	fmt.Printf("#   with parameters: %v\n", strings.Join(args, " "))
}

func CheckApiLevel(section *ini.Section, min int, max int) {
	if !section.HasKey(ApiLevelKey) {
		log.Println("ERROR: Unable to check compatibility, environment does not specify API level.")
		os.Exit(1)
	}

	apiLevelString := section.Key(ApiLevelKey).Value()
	apiLevel, err := strconv.Atoi(apiLevelString)
	if err != nil {
		log.Printf("ERROR: Unable to check compatibility, as API level %v could not be read as integer.\n", apiLevelString)
		os.Exit(1)
	}

	supportedInterval := ""
	if min != 0 {
		supportedInterval += fmt.Sprintf("%v <= ", min)
	}
	supportedInterval += "level"
	if max != 0 {
		supportedInterval += fmt.Sprintf(" <= %v", max)
	}

	if (min != 0 && apiLevel < min) || (max != 0 && apiLevel > max) {
		log.Printf("ERROR: API level %v is not within the supported interval for this command: %v\n", apiLevel, supportedInterval)
		os.Exit(1)
	}
}
