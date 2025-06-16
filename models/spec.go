// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package models

type Spec struct {
	Path string `json:"path"`

	// other fields... :)
}

func (s Spec) GetPath() string {
	return s.Path
}
