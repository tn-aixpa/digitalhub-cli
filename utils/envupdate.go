// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
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
	config, err := FetchConfig(baseEndpoint + "/.well-known/configuration")
	if err != nil {
		return
	}
	for k, v := range config {
		newKey := k
		if newKey == ClientIdKey {
			newKey = "client_id"
		}
		UpdateKey(section, newKey, v)
	}

	// OpenID Configuration
	openIdConfig, err := FetchConfig(baseEndpoint + "/.well-known/openid-configuration")
	if err != nil {
		return
	}
	for _, k := range OpenIdFields {
		if v, ok := openIdConfig[k]; ok && v != "" {
			UpdateKey(section, k, v)
		}
	}

	// Update timestamp
	section.Key(UpdatedEnvKey).SetValue(time.Now().Format(time.RFC3339))
	SaveIni(cfg)
}

func UpdateKey(section *ini.Section, k string, v interface{}) {
	if !section.HasKey(k) {
		section.NewKey(k, ReflectValue(v))
	} else {
		section.Key(k).SetValue(ReflectValue(v))
	}
}
