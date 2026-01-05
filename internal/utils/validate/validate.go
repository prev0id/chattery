package validate

import (
	"net/mail"

	"chattery/internal/utils/errors"
)

const (
	UsernameFieldName = "username"
	PasswordFieldName = "password"
	LoginFieldName    = "login"
)

func Username(username string) error {
	if err := minLength(username, 5, UsernameFieldName); err != nil {
		return err
	}
	if err := maxLength(username, 20, UsernameFieldName); err != nil {
		return err
	}
	if err := startWithLowercaseLetter(username, UsernameFieldName); err != nil {
		return err
	}
	if err := endWithLowercaseLetter(username, UsernameFieldName); err != nil {
		return err
	}
	if err := containsOnlyLowerCaseAndUnderscore(username, UsernameFieldName); err != nil {
		return err
	}
	return nil
}

func Login(login string) error {
	if err := validEmail(login, LoginFieldName); err != nil {
		return err
	}
	return nil
}

func Password(password string) error {
	if err := minLength(password, 8, PasswordFieldName); err != nil {
		return err
	}
	if err := maxLength(password, 32, PasswordFieldName); err != nil {
		return err
	}
	if err := hasLowerCaseLetter(password, PasswordFieldName); err != nil {
		return err
	}
	if err := hasUpperCaseLetter(password, PasswordFieldName); err != nil {
		return err
	}
	if err := hasDigit(password, PasswordFieldName); err != nil {
		return err
	}
	return nil
}

func minLength(str string, min int, field string) error {
	if len(str) < min {
		return errors.E().
			Kind(errors.InvalidRequest).
			Messagef("%s must be at least %d characters long", field, min)
	}
	return nil
}

func maxLength(str string, max int, field string) error {
	if len(str) > max {
		return errors.E().
			Kind(errors.InvalidRequest).
			Messagef("%s must be at most %d characters long", field, max)
	}
	return nil
}

func startWithLowercaseLetter(str string, field string) error {
	if str[0] < 'a' || str[0] > 'z' {
		return errors.E().
			Kind(errors.InvalidRequest).
			Messagef("%s must start with a lowercase letter", field)
	}
	return nil
}

func endWithLowercaseLetter(str, field string) error {
	lastIdx := len(str) - 1
	if str[lastIdx] < 'a' || str[lastIdx] > 'z' {
		return errors.E().
			Kind(errors.InvalidRequest).
			Messagef("%s must end with a lowercase letter", field)
	}
	return nil
}

func hasLowerCaseLetter(str, field string) error {
	for _, c := range str {
		if 'a' <= c && c <= 'z' {
			return nil
		}
	}
	return errors.E().
		Kind(errors.InvalidRequest).
		Messagef("%s must contain at least one lowercase letter", field)
}

func hasUpperCaseLetter(str, field string) error {
	for _, c := range str {
		if 'A' <= c && c <= 'Z' {
			return nil
		}
	}
	return errors.E().
		Kind(errors.InvalidRequest).
		Messagef("%s must contain at least one uppercase letter", field)
}

func hasDigit(str, field string) error {
	for _, c := range str {
		if '0' <= c && c <= '9' {
			return nil
		}
	}
	return errors.E().
		Kind(errors.InvalidRequest).
		Messagef("%s must contain at least one digit", field)
}

func validEmail(str, field string) error {
	_, err := mail.ParseAddress(str)
	if err == nil {
		return nil
	}
	return errors.E(err).
		Kind(errors.InvalidRequest).
		Messagef("%s must be a valid email address", LoginFieldName)
}

func containsOnlyLowerCaseAndUnderscore(str, field string) error {
	invalid := false
	for _, char := range str {
		if char != '_' && (char < 'a' || char > 'z') {
			invalid = true
			break
		}
	}
	if !invalid {
		return nil
	}
	return errors.E().
		Kind(errors.InvalidRequest).
		Messagef("%s can only contain lowercase letters (a-z) and underscores", field)
}

func NotEmpty[T comparable](value T, field string) error {
	var empty T
	if value != empty {
		return nil
	}
	return errors.E().
		Kind(errors.InvalidRequest).
		Messagef("%s must be provided", field)
}
