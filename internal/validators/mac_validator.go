package validators

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// Constant defining the expected size of a MAC address in bytes.
const macBytes = 6 // macBytes is the standard length of a 48-bit MAC address.

// ValidateMAC validates a MAC address string for correctness.
// It trims whitespace, ensures the address is non-empty, and parses it into a 48-bit (6-byte) MAC address.
// Returns an error if the MAC address is invalid or not exactly 6 bytes.
func ValidateMAC(mac string) error {
	mac = strings.TrimSpace(mac)
	if mac == "" {
		return errors.New("MAC address is required")
	}

	macAddr, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("parsing MAC address: %w", err)
	}
	if len(macAddr) != macBytes {
		return fmt.Errorf("MAC address must be %d bytes, got %d", macBytes, len(macAddr))
	}
	return nil
}
