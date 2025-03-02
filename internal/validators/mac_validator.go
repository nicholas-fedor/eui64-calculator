package validators

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// Constant defining the maximum string length for a MAC address.
const (
	macStrLen = 17 // macStrLen is the maximum string length for "xx-xx-xx-xx-xx-xx".
)

// ValidateMAC validates a MAC address string for correctness.
// It trims whitespace, ensures the address is non-empty, checks the string length,
// and parses it into a valid MAC address using net.ParseMAC.
// Returns an error if the MAC address is invalid or exceeds the maximum length.
func ValidateMAC(macStr string) error {
	macStr = strings.TrimSpace(macStr)
	if macStr == "" {
		return errors.New("MAC address is required")
	}

	if len(macStr) > macStrLen {
		return fmt.Errorf("MAC address string exceeds maximum length of %d characters", macStrLen)
	}

	_, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("parsing MAC address: %w", err)
	}

	return nil
}
