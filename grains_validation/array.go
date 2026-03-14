package grains_validation

// ValidateArrayIntegers checks if all integers in the slice are within [minValue, maxValue].
// Returns false if the slice is empty or any value is out of range.
func ValidateArrayIntegers(values []int, minValue, maxValue int) bool {
	if len(values) == 0 {
		return false // must contain at least one integer
	}

	for _, n := range values {
		if n < minValue || n > maxValue {
			return false // out of allowed range
		}
	}

	return true
}
