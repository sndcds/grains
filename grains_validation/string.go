package grains_validation

import (
	"strconv"
	"strings"
)

// ValidateCategories checks if a comma-separated string contains valid integers
// within the range [minValue, maxValue]. Returns true if valid.
func ValidateStringIntegers(value, separator string, minValue, maxValue int) bool {
	if value == "" {
		return false
	}

	parts := strings.Split(value, separator)
	for _, p := range parts {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil {
			return false // not a valid integer
		}
		if n < minValue || n > maxValue {
			return false // out of allowed range
		}
	}

	return true
}
