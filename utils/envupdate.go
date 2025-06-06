package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gopkg.in/ini.v1"
)

func CheckUpdateEnvironment(cfg *ini.File, section *ini.Section) {
	if section.HasKey(UpdatedEnvKey) {
		updated, err := time.Parse(time.RFC3339, section.Key(UpdatedEnvKey).Value())
		if err != nil || updated.Add(outdatedAfterHours*time.Hour).Before(time.Now()) {
			updateEnvironment(cfg, section)
		}
	}
}

func updateEnvironment(cfg *ini.File, section *ini.Section) {
	baseEndpoint := section.Key(DhCoreEndpoint).Value()
	if baseEndpoint == "" {
		return
	}

	// Configuration
	config, err := doGet(baseEndpoint + "/.well-known/configuration")
	if err != nil {
		return
	}
	updateKeys(section, config)

	// OpenID Configuration
	openIdConfig, err := doGet(baseEndpoint + "/.well-known/openid-configuration")
	if err != nil {
		return
	}
	updateKeys(section, openIdConfig)

	// Update timestamp
	section.Key(UpdatedEnvKey).SetValue(time.Now().Format(time.RFC3339))
	SaveIni(cfg)
}

func doGet(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func updateKeys(section *ini.Section, config map[string]interface{}) {
	for k, v := range config {
		if !section.HasKey(k) {
			section.NewKey(k, ReflectValue(v))
		} else {
			section.Key(k).SetValue(ReflectValue(v))
		}
	}
}
