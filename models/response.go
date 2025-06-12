// SPDX-FileCopyrightText: © 2025 DSLab - Fondazione Bruno Kessler
//
// SPDX-License-Identifier: Apache-2.0

package models

type Response[T BaseModel] struct {
	Content []T `json:"content"`
}
