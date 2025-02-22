package validators

import (
	"fmt"
	"net"
	"strings"
	"unicode"
)

const (
	MacLength        = 6
	MaxPrefixHextets = 4
	MaxHextetLength  = 4
)

// CombinedValidator implements validation for MAC addresses and IPv6 prefixes.
type CombinedValidator struct{}

// ValidateMAC checks if the MAC address is valid.
func (v *CombinedValidator) ValidateMAC(mac string) error {
	mac = strings.TrimSpace(mac)
	if mac == "" {
		return ErrMACRequired
	}

	macAddr, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrParseMACFailed, err)
	}

	if len(macAddr) != MacLength {
		return fmt.Errorf("%w: got %d bytes", ErrInvalidMACLength, len(macAddr))
	}

	return nil
}

// ValidateIPv6Prefix validates an IPv6 prefix.
func (v *CombinedValidator) ValidateIPv6Prefix(prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ErrPrefixRequired
	}

	prefixParts := strings.Split(strings.TrimSuffix(prefix, "::"), ":")
	if len(prefixParts) > MaxPrefixHextets {
		return fmt.Errorf("%w: got %d parts", ErrPrefixTooLong, len(prefixParts))
	}

	return validateHextets(prefixParts, prefix)
}

// validateHextets checks each hextet for validity.
func validateHextets(prefixParts []string, _ string) error {
	for i, hextet := range prefixParts {
		if hextet == "" {
			if i != 0 && i != len(prefixParts)-1 {
				return fmt.Errorf("%w: %v", ErrEmptyHextet, prefixParts)
			}

			continue
		}

		if len(hextet) > MaxHextetLength {
			return fmt.Errorf("%w: %v", ErrInvalidHextetLength, hextet)
		}

		if err := isValidHextetChar(hextet); err != nil {
			return err
		}
	}

	return nil
}

// isValidHextetChar ensures all characters in a hextet are valid hex digits.
func isValidHextetChar(hextet string) error {
	for _, char := range hextet {
		if !unicode.IsDigit(char) && (char < 'A' || char > 'F') && (char < 'a' || char > 'f') {
			return fmt.Errorf("%w: %v", ErrInvalidHextetChar, hextet)
		}
	}

	return nil
}
