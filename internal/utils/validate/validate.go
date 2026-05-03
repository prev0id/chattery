package validate

import (
	"net/mail"
	"regexp"
	"slices"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

const (
	UsernameFieldName = "username"
	NameFieldName     = "name"
	TypeFieldName     = "type"
	PasswordFieldName = "password"
	LoginFieldName    = "login"
)

var (
	onlyWords = regexp.MustCompile(`^[\w- ]+$`)
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

	return containsOnlyLowerCaseAndUnderscore(username, UsernameFieldName)
}

func Login(login string) error {
	return validEmail(login, LoginFieldName)
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

	return hasDigit(password, PasswordFieldName)
}

func ServerName(name string) error {
	if err := minLength(name, 5, NameFieldName); err != nil {
		return err
	}
	if err := maxLength(name, 25, NameFieldName); err != nil {
		return err
	}
	return containsOnlyWords(name, NameFieldName)
}

func TopicName(name string) error {
	if err := minLength(name, 2, NameFieldName); err != nil {
		return err
	}
	if err := maxLength(name, 20, NameFieldName); err != nil {
		return err
	}
	return containsOnlyWords(name, NameFieldName)
}

func TopicType(topic string) error {
	return oneOf(
		topic,
		TypeFieldName,
		domain.TopicTypeText.String(),
		domain.TopicTypeVoice.String(),
	)
}

func minLength(str string, length int, field string) error {
	if len(str) < length {
		return errutil.E().
			Kind(errutil.InvalidRequest).
			Messagef("%s must be at least %d characters long", field, length)
	}
	return nil
}

func maxLength(str string, length int, field string) error {
	if len(str) > length {
		return errutil.E().
			Kind(errutil.InvalidRequest).
			Messagef("%s must be at most %d characters long", field, length)
	}
	return nil
}

func startWithLowercaseLetter(str string, field string) error {
	if str[0] < 'a' || str[0] > 'z' {
		return errutil.E().
			Kind(errutil.InvalidRequest).
			Messagef("%s must start with a lowercase letter", field)
	}
	return nil
}

func endWithLowercaseLetter(str, field string) error {
	lastIdx := len(str) - 1
	if str[lastIdx] < 'a' || str[lastIdx] > 'z' {
		return errutil.E().
			Kind(errutil.InvalidRequest).
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
	return errutil.E().
		Kind(errutil.InvalidRequest).
		Messagef("%s must contain at least one lowercase letter", field)
}

func hasUpperCaseLetter(str, field string) error {
	for _, c := range str {
		if 'A' <= c && c <= 'Z' {
			return nil
		}
	}
	return errutil.E().
		Kind(errutil.InvalidRequest).
		Messagef("%s must contain at least one uppercase letter", field)
}

func hasDigit(str, field string) error {
	for _, c := range str {
		if '0' <= c && c <= '9' {
			return nil
		}
	}
	return errutil.E().
		Kind(errutil.InvalidRequest).
		Messagef("%s must contain at least one digit", field)
}

func validEmail(str, field string) error {
	_, err := mail.ParseAddress(str)
	if err == nil {
		return nil
	}
	return errutil.E(err).
		Kind(errutil.InvalidRequest).
		Messagef("%s must be a valid email address", field)
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
	return errutil.E().
		Kind(errutil.InvalidRequest).
		Messagef("%s can only contain lowercase letters (a-z) and underscores", field)
}

func containsOnlyWords(str, field string) error {
	if onlyWords.MatchString(str) {
		return nil
	}
	return errutil.E().
		Kind(errutil.InvalidRequest).
		Messagef("%s can only contain letters (a-z, A-Z), digits, spaces, underscores and dashes -", field)
}

func oneOf[T comparable](value T, field string, targets ...T) error {
	if slices.Contains(targets, value) {
		return nil
	}

	return errutil.E().
		Kind(errutil.InvalidRequest).
		Messagef("%s must be one of values %v", field, targets)
}

func NotEmpty[T comparable](value T, field string) error {
	var empty T
	if value != empty {
		return nil
	}
	return errutil.E().
		Kind(errutil.InvalidRequest).
		Messagef("%s must be provided", field)
}
