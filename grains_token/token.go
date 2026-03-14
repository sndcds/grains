package grains_token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"math/big"
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

// encodeBase62 converts a byte slice into a Base62 string
func encodeBase62Bytes(b []byte) string {
	n := new(big.Int).SetBytes(b) // interpret bytes as a big integer
	if n.Sign() == 0 {
		return "0"
	}

	result := ""
	base := big.NewInt(62)
	mod := new(big.Int)

	for n.Sign() > 0 {
		n.DivMod(n, base, mod)
		result = string(base62Chars[mod.Int64()]) + result
	}
	return result
}

// RandomishToken generates a deterministic, unique, random-looking token
// from a timestamp and a unique row ID, optionally adding a prefix or suffix.
// The result is URL-friendly.
func RandomishToken(t time.Time, id int64, prefix string, suffix string) string {
	epochUs := uint64(t.UnixNano() / 1000)
	combined := epochUs ^ (uint64(id) * 37)
	combined = (combined << 23) | (combined >> (64 - 23))
	token := encodeBase62(combined)
	return prefix + token + suffix
}

func SecureToken(t time.Time, id int64, secret []byte, prefix, suffix string) string {
	msg := fmt.Sprintf("%d:%d", t.UnixNano(), id)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(msg))
	sum := mac.Sum(nil)
	token := base64.RawURLEncoding.EncodeToString(sum[:24])
	return prefix + token + suffix
}

func GenerateShortAPIKey(id int64, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(fmt.Sprintf("%d", id)))
	sum := mac.Sum(nil)
	// Convert first 8 bytes to uint64
	n := binary.BigEndian.Uint64(sum[:8])
	// Encode in Base62
	return encodeBase62(n)
}

// GenerateUniqueAPIKey creates a unique, cryptographically strong API key for a given id
func GenerateUniqueAPIKey(id int64, secret []byte) string {
	// Step 1: Create HMAC-SHA256 of the ID
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(fmt.Sprintf("%d", id)))
	sum := mac.Sum(nil) // 32 bytes

	// Step 2: Encode full 32 bytes in Base62
	return encodeBase62Bytes(sum)
}

func GenerateUuid() string {
	id := uuid.New()
	return encodeBase62Bytes(id[:])
}
