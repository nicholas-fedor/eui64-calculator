package validators

import (
	"errors"
	"fmt"
	"strings"
)

// Constants defining constraints for IPv6 prefix validation.
const (
	maxHextets      = 4 // maxHextets is the maximum number of hextets allowed in an IPv6 prefix.
	maxHextetLength = 4 // maxHextetLength is the maximum length of a single hextet in characters.
)

// ValidateIPv6Prefix validates an IPv6 prefix string for correctness.
// It trims whitespace, ensures the prefix is non-empty, and checks that it contains 4 or fewer hextets,
// with each hextet being valid hexadecimal (up to 4 characters) and no internal empty hextets.
// Allows trailing "::" for zero compression. Returns an error if the prefix is invalid.
func ValidateIPv6Prefix(prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("a non-blank IPv6 prefix is expected")
	}

	prefix = strings.TrimSuffix(prefix, "::")
	if prefix == "" {
		return nil // "::" alone is valid, implying all zeros.
	}

	hextets := strings.Split(prefix, ":")
	if len(hextets) > maxHextets {
		return fmt.Errorf("IPv6 prefix must be %d or fewer hextets", maxHextets)
	}

	for i, hextet := range hextets {
		if hextet == "" {
			if i != 0 && i != len(hextets)-1 {
				return errors.New("empty hextet in IPv6 prefix")
			}
			continue
		}
		for _, char := range hextet {
			if !isHexDigit(char) {
				return errors.New("invalid character in hextet")
			}
		}
		if len(hextet) > maxHextetLength {
			return errors.New("invalid hextet length in IPv6 prefix")
		}
	}
	return nil
}

// isHexDigit reports whether a rune is a valid hexadecimal digit.
// It checks if the character is 0-9, a-f, or A-F, returning true if valid, false otherwise.
func isHexDigit(char rune) bool {
	return (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')
}
