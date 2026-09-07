package grains_validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nbutton23/zxcvbn-go"
)

const (
	// MaxPasswordBytes is the maximum password size supported by bcrypt.
	//
	// bcrypt operates on at most 72 bytes. This is a byte limit, not
	// a Unicode character/rune limit.
	MaxPasswordBytes = 72

	// DefaultMinPasswordLength is the minimum password length required
	// by the application.
	DefaultMinPasswordLength = 12
)

var (
	ErrPasswordTooShort      = errors.New("password is too short")
	ErrPasswordTooLong       = errors.New("password is too long")
	ErrPasswordTooWeak       = errors.New("password is too weak")
	ErrPasswordContainsEmail = errors.New(
		"password contains the email name",
	)
	ErrPasswordCommon = errors.New("password is too common")
)

var commonPasswords = map[string]struct{}{
	"password":    {},
	"123456":      {},
	"123456789":   {},
	"qwerty":      {},
	"password123": {},
	"admin":       {},
	"letmein":     {},
}

// ValidatePassword validates a password according to the application's
// password policy.
//
// The password itself is never modified. In particular, whitespace,
// Unicode characters and casing are preserved exactly as supplied by
// the user.
//
// The email is only used to prevent the email's local part from being
// included in the password.
func ValidatePassword(email, password string, minLength int) error {
	// Password length is deliberately measured in bytes because bcrypt's
	// maximum input size is 72 bytes.
	if utf8.RuneCountInString(password) < minLength {
		return fmt.Errorf(
			"%w: minimum length is %d characters",
			ErrPasswordTooShort,
			minLength,
		)
	}

	if len(password) > MaxPasswordBytes {
		return fmt.Errorf(
			"%w: maximum length is %d bytes",
			ErrPasswordTooLong,
			MaxPasswordBytes,
		)
	}

	// zxcvbn evaluates the password without modifying it.
	result := zxcvbn.PasswordStrength(password, nil)
	if result.Score < 3 {
		return ErrPasswordTooWeak
	}

	// Do not reject very short email local parts such as "a" or "i".
	// They would occur naturally in almost every password.
	if email != "" {
		localPart := strings.ToLower(strings.SplitN(email, "@", 2)[0])

		if len(localPart) >= 3 &&
			strings.Contains(strings.ToLower(password), localPart) {
			return ErrPasswordContainsEmail
		}
	}

	if isCommonPassword(password) {
		return ErrPasswordCommon
	}

	return nil
}

func isCommonPassword(password string) bool {
	_, ok := commonPasswords[strings.ToLower(password)]
	return ok
}
