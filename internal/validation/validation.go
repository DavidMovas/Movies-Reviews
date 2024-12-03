package validation

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/DavidMovas/Movies-Reviews/contracts"

	"gopkg.in/validator.v2"
)

var (
	passwordMinLength       = 8
	emailMaxLength          = 127
	passwordSpecialChars    = "!#$%&'*+/=?^_`{|}~@"
	passwordRequiredEntries = []struct {
		name  string
		chars string
	}{
		{"lowercase", "abcdefghijklmnopqrstuvwxyz"},
		{"uppercase", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"numbers", "0123456789"},
		{"special character (" + passwordSpecialChars + ")", passwordSpecialChars},
	}
)

func SetupValidators() {
	validators := []struct {
		name string
		fn   validator.ValidationFunc
	}{
		{"email", email},
		{"password", password},
		{"nonzero", nonzero},
		{"role", role},
	}

	for _, v := range validators {
		_ = validator.SetValidationFunc(v.name, v.fn)
	}
}

func password(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("password must be a string")
	}

	if len(s) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters long", passwordMinLength)
	}

	for _, entry := range passwordRequiredEntries {
		if !strings.ContainsAny(s, entry.chars) {
			return fmt.Errorf("password must contain at least one of the following required entries: %s", entry.name)
		}
	}

	return nil
}

func email(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("email must be a string")
	}

	if len(s) > emailMaxLength {
		return fmt.Errorf("email must be at most %d characters long", emailMaxLength)
	}

	_, err := mail.ParseAddress(s)
	return err
}

func nonzero(v interface{}, _ string) error {
	s, ok := v.(string)
	if ok && s == "" {
		return fmt.Errorf("must not be empty")
	}
	return nil
}

func role(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("role must be a string")
	}

	if !contracts.ValidateRole(s) {
		return fmt.Errorf("invalid role")
	}
	return nil
}
