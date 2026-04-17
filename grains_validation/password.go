package grains_validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nbutton23/zxcvbn-go"
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

func ValidatePassword(email, password string, minLength int) error {
	// 1. Basic length check (fast fail)
	if len(password) < minLength {
		return fmt.Errorf("password must be at least %d characters long", minLength)
	}

	// 2. Strength check (zxcvbn score 0–4)
	result := zxcvbn.PasswordStrength(password, nil)
	if result.Score < 3 {
		return errors.New("password is too weak")
	}

	// 3. Avoid using email name inside password
	if email != "" {
		localPart := strings.ToLower(strings.Split(email, "@")[0])
		if strings.Contains(strings.ToLower(password), localPart) {
			return errors.New("password should not contain your email name")
		}
	}

	// 4. Common password quick reject (optional but useful)
	if isCommonPassword(password) {
		return errors.New("password is too common")
	}

	// 5. Optional: breached password check hook
	// if isPwned(password) {
	//     return errors.New("password has appeared in data breaches")
	// }

	return nil
}

func isCommonPassword(pw string) bool {
	_, ok := commonPasswords[strings.ToLower(pw)]
	return ok
}
