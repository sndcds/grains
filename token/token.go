package token

import (
	"fmt"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// encodeBase62 converts a uint64 to a base62 string
func encodeBase62(n uint64) string {
	if n == 0 {
		return "0"
	}

	result := ""
	for n > 0 {
		remainder := n % 62
		result = string(base62Chars[remainder]) + result
		n /= 62
	}
	return result
}

// RandomishToken generates a deterministic, unique, random-looking token
// from a timestamp and a unique row ID, optionally adding a prefix or suffix.
// The result is URL-friendly.
func RandomishToken(t time.Time, rowID int64, prefix string, suffix string) string {
	// Step 1: Convert timestamp to microseconds since epoch
	epochUs := uint64(t.UnixNano() / 1000)

	// Step 2: Combine with rowID using XOR (multiplier spreads bits)
	combined := epochUs ^ (uint64(rowID) * 37)

	// Step 3: Bit rotation left by 23 for diffusion
	combined = (combined << 23) | (combined >> (64 - 23))

	// Step 4: Encode to base62
	token := encodeBase62(combined)

	// Step 5: Add optional prefix/suffix
	return prefix + token + suffix
}

func main() {
	t := time.Now()
	rowID := int64(42)

	// Example 1: token without prefix/suffix
	token1 := RandomishToken(t, rowID, "", "")
	fmt.Println("Token1:", token1)

	// Example 2: token with prefix
	token2 := RandomishToken(t, rowID, "img_", "")
	fmt.Println("Token2:", token2)

	// Example 3: token with suffix
	token3 := RandomishToken(t, rowID, "", "_2026")
	fmt.Println("Token3:", token3)

	// Example 4: token with prefix and suffix
	token4 := RandomishToken(t, rowID, "img_", "_2026")
	fmt.Println("Token4:", token4)
}
