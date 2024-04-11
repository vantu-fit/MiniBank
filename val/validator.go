package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\S]+$`).MatchString
)

func ValidateString(value string, minLenght int, maxLenght int) error {
	n := len(value)
	if n < minLenght || n > maxLenght {
		return fmt.Errorf("must contain from %d-%d characters", minLenght, maxLenght)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only  lowcase letters , digits or underscore")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil

}

func ValidateFullname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil

}

func ValidateEmailId(value int64) error {
	if value < 0 {
		return fmt.Errorf("must be s positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 30, 128)
}
