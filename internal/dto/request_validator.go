package dto

import (
	"errors"
	"reflect"
	"strconv"
)

// Request is the interface all request types should satisfy when validation is needed.
type Request interface {
	Validate() error
}

// ValidateRequired returns an error if the value is nil or an empty string.
func ValidateRequired(value interface{}, fieldName string) error {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.String && value == "") {
		return errors.New(fieldName + " is required")
	}
	return nil
}

// ValidateEmail returns an error if the string does not look like a valid e-mail.
func ValidateEmail(value string) error {
	if value == "" {
		return nil
	}
	if len(value) < 3 || !contains(value, "@") || !contains(value, ".") {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateMinLength returns an error if the string is shorter than min characters.
func ValidateMinLength(value string, min int) error {
	if len(value) < min {
		return errors.New("must be at least " + strconv.Itoa(min) + " characters")
	}
	return nil
}

// ValidateLoginRequest validates the fields of a LoginRequest.
func ValidateLoginRequest(r LoginRequest) error {
	if err := ValidateRequired(r.Email, "Email"); err != nil {
		return err
	}
	if err := ValidateEmail(r.Email); err != nil {
		return err
	}
	if err := ValidateRequired(r.Password, "Password"); err != nil {
		return err
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

// ValidateRegisterRequest validates the fields of a RegisterRequest.
func ValidateRegisterRequest(r RegisterRequest) error {
	if err := ValidateRequired(r.Name, "Name"); err != nil {
		return err
	}
	if len(r.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if err := ValidateRequired(r.Email, "Email"); err != nil {
		return err
	}
	if err := ValidateEmail(r.Email); err != nil {
		return err
	}
	if err := ValidateRequired(r.Password, "Password"); err != nil {
		return err
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if err := ValidateRequired(r.ConfirmPassword, "Confirm Password"); err != nil {
		return err
	}
	if r.Password != r.ConfirmPassword {
		return errors.New("password and confirm password do not match")
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
