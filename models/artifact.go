// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package models

type Artifact struct {
	Spec Spec `json:"spec"`
}

func (a Artifact) GetSpec() Spec {
	return a.Spec
}
