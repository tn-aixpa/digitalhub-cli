package utils

import "errors"

// Checks if value is empty and returns error if it is
func ValidateRequired(value, fieldName string) error {
	if value == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}
