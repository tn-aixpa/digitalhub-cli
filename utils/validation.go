package utils

import "errors"

// ValidateRequired controlla se value e' empty nel caso ritorna errore
func ValidateRequired(value, fieldName string) error {
	if value == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}
