package cmd

import (
	"bufio"
	"dhcli/utils"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"gopkg.in/ini.v1"
)

type OpenIDConfig struct {
	AuthorizationEndpoint string `json:"authorization_endpoint" ini:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint" ini:"token_endpoint"`
	Issuer                string `json:"issuer" ini:"issuer"`
	ClientID              string `json:"dhcore_client_id" ini:"dhcore_client_id"`
	Scope                 string `json:"scope" ini:"scope"`
	AccessToken           string `json:"access_token" ini:"access_token"`
	RefreshToken          string `json:"refresh_token" ini:"refresh_token"`
}

type CoreConfig struct {
	Name     string `json:"dhcore_name" ini:"dhcore_name"`
	Issuer   string `json:"issuer" ini:"issuer"`
	Version  string `json:"dhcore_version" ini:"dhcore_version"`
	ClientID string `json:"dhcore_client_id" ini:"dhcore_client_id"`
}

func init() {
	RegisterCommand(&Command{
		Name:        "register",
		Description: "./dhcli register [-n <name>] <endpoint>",
		SetupFlags: func(fs *flag.FlagSet) {
			fs.String("n", "", "name")
		},
		Handler: registerHandler,
	})
}

func registerHandler(args []string, fs *flag.FlagSet) {
	ini.DefaultHeader = true

	if len(args) < 1 {
		log.Fatalf("Error: Endpoint is required.\nUsage: ./dhcli register [-n <name>] <endpoint>")
	}
	fs.Parse(args)

	name := fs.Lookup("n").Value.String()
	endpoint := fs.Args()[0]

	// Read or initialize ini file
	cfg := utils.LoadIni(true)

	//collect to map+struct
	res, coreConfig := fetchConfig(endpoint + "/.well-known/configuration")
	if name == "" || name == "null" {
		name = coreConfig.Name
		if name == "" {
			log.Fatalf("Failed to register: environment name not specified and not defined in core's configuration.")
		}
	}
	sec := cfg.Section(name)
	sec.ReflectFrom(&coreConfig)

	// Fetch OpenID configuration
	openIDConfig := fetchOpenIDConfig(endpoint + "/.well-known/openid-configuration")
	openIDConfig.ClientID = coreConfig.ClientID
	sec.ReflectFrom(&openIDConfig)

	for k, v := range res {
		//add missing keys
		if !sec.HasKey(k) {
			f := reflect.ValueOf(v)
			var val string
			switch f.Kind() {
			case reflect.String:
				val = f.String()
			case reflect.Int, reflect.Int64:
				val = fmt.Sprint(f.Int())
			case reflect.Uint, reflect.Uint64:
				val = fmt.Sprint(f.Uint())
			case reflect.Float64:
				val = fmt.Sprint(f.Float())
			case reflect.Bool:
				val = fmt.Sprint(f.Bool())
			case reflect.TypeOf(time.Now()).Kind():
				val = f.Interface().(time.Time).Format(time.RFC3339)
			case reflect.Slice:
				val = fmt.Sprint(f.Interface())
			default:
				val = ""
			}

			sec.NewKey(k, val)
		}
	}

	//check for default env
	dsec := cfg.Section("DEFAULT")
	if !dsec.HasKey(utils.CurrentEnvironment) {
		dsec.NewKey(utils.CurrentEnvironment, name)
	}

	// gitignoreAddIniFile()
	utils.SaveIni(cfg)
	log.Printf("'%v' registered.", name)
}

func fetchConfig(configURL string) (map[string]interface{}, CoreConfig) {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Fatalf("Error fetching core configuration: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Core responded with error %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading core configuration response: %v", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Fatalf("Error parsing core configuration: %v", err)
	}

	var config CoreConfig
	if err := json.Unmarshal(body, &config); err != nil {
		log.Fatalf("Error parsing core configuration: %v", err)
	}

	return res, config
}

func fetchOpenIDConfig(configURL string) OpenIDConfig {
	resp, err := http.Get(configURL)
	if err != nil {
		log.Fatalf("Error fetching OpenID configuration: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Core responded with error %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading OpenID configuration response: %v", err)
	}

	var config OpenIDConfig
	if err := json.Unmarshal(body, &config); err != nil {
		log.Fatalf("Error parsing OpenID configuration: %v", err)
	}

	return config
}

func toMap(strc interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})

	// get or dereference
	val := reflect.ValueOf(strc)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	if val.Kind() != reflect.Struct {
		return res, errors.New("variable given is not a struct or a pointer to a struct")
	}

	//export to value
	//NOTE: doesn't support nested structs
	for i := 0; i < val.NumField(); i++ {
		fName := typ.Field(i).Name
		fValue := val.Field(i).Interface()
		res[fName] = fValue
	}

	return res, nil
}

func gitignoreAddIniFile() {
	path := "./.gitignore"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Cannot open .gitignore file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == utils.IniName {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading .gitignore file contents: %v", err)
	}

	if _, err = f.WriteString(utils.IniName); err != nil {
		log.Fatalf("Error while adding entry to .gitignore file: %v", err)
	}
}
