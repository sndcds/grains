package grains_uuid

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"
)

// Standard UUID regex (36 chars with hyphens)
var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func Uuidv7() ([16]byte, error) {
	// Random bytes
	var value [16]byte
	_, err := rand.Read(value[:])
	if err != nil {
		return value, err
	}

	// Current timestamp in ms
	timestamp := big.NewInt(time.Now().UnixMilli())

	// Timestamp
	timestamp.FillBytes(value[0:6])

	// Version and variant
	value[6] = (value[6] & 0x0F) | 0x70
	value[8] = (value[8] & 0x3F) | 0x80

	return value, nil
}

func Uuidv7String() (string, error) {
	uuid, err := Uuidv7()
	if err != nil {
		return "", err
	}
	return Uuidv7ToString(uuid), nil
}

func Uuidv7ToString(u [16]byte) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4],
		u[4:6],
		u[6:8],
		u[8:10],
		u[10:16],
	)
}

// Validate and decode standard UUIDv7 string
func Uuidv7FromString(s string) ([16]byte, error) {
	var u [16]byte

	if uuidRegex.MatchString(s) {
		// Remove hyphens
		hexStr := s[0:8] + s[9:13] + s[14:18] + s[19:23] + s[24:36]
		b, err := hex.DecodeString(hexStr)
		if err != nil {
			return u, err
		}
		copy(u[:], b)
		return u, nil
	}

	// Try Base64 URL-safe variant
	if len(s) == 22 { // 16 bytes in URL-safe base64 without padding
		b, err := base64.RawURLEncoding.DecodeString(s)
		if err != nil {
			return u, err
		}
		if len(b) != 16 {
			return u, errors.New("invalid UUID length after base64 decode")
		}
		copy(u[:], b)
		return u, nil
	}

	return u, errors.New("invalid UUIDv7 string")
}
