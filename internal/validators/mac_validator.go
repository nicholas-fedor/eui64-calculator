// Package validators provides validation functions for input data in the EUI-64 calculator application.
// It includes validators for MAC addresses and IPv6 prefixes, ensuring they meet format and length requirements
// before proceeding with calculations.
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

// Static error variables.
var (
	ErrMACRequired      = errors.New("MAC address is required")
	ErrMACLengthExceeds = fmt.Errorf(
		"MAC address string exceeds maximum length of %d characters",
		macStrLen,
	)
	ErrMACParseFailed = errors.New("parsing MAC address")
)

// ValidateMAC validates a MAC address string for correctness.
// It trims whitespace, ensures the address is non-empty, checks the string length,
// and parses it into a valid MAC address using net.ParseMAC.
// Returns an error if the MAC address is invalid or exceeds the maximum length.
func ValidateMAC(macStr string) error {
	macStr = strings.TrimSpace(macStr)
	if macStr == "" {
		return ErrMACRequired
	}

	if len(macStr) > macStrLen {
		return ErrMACLengthExceeds
	}

	_, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMACParseFailed, err)
	}

	return nil
}
