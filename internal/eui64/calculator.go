package eui64

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	ipv6Hextets      = 8
	PrefixMaxHextets = 4
	macBytes         = 6
	eui64Bytes       = 8
	fffeMarkerLow    = 0xFF
	fffeMarkerHigh   = 0xFE
	byteShift        = 8
	hextetIndex4     = 4
	hextetIndex5     = 5
	hextetIndex6     = 6
	hextetIndex7     = 7
)

var (
	ErrParseMACFailed     = errors.New("parsing MAC address failed")
	ErrInvalidMACLength   = errors.New("MAC address must be 6 bytes")
	ErrPrefixTooLong      = errors.New("IPv6 prefix exceeds 4 hextets")
	ErrInvalidEmptyHextet = errors.New("invalid empty hextet in IPv6 prefix")
	ErrInvalidHextet      = errors.New("invalid hextet in IPv6 prefix")
)

// Calculator defines the interface for computing EUI-64 identifiers and IPv6 addresses.
type Calculator interface {
	// CalculateEUI64 computes the EUI-64 interface ID and full IPv6 address from a MAC address and IPv6 prefix.
	CalculateEUI64(mac, prefix string) (string, string, error)
}

// DefaultCalculator implements the Calculator interface using the standard EUI-64 algorithm.
type DefaultCalculator struct{}

// CalculateEUI64 computes the EUI-64 interface ID and full IPv6 address from a MAC address and prefix.
// It delegates to the standalone CalculateEUI64 function.
func (d *DefaultCalculator) CalculateEUI64(mac, prefix string) (string, string, error) {
	return CalculateEUI64(mac, prefix)
}

// CalculateEUI64 computes the EUI-64 interface ID and full IPv6 address.
func CalculateEUI64(macStr, prefixStr string) (string, string, error) {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return "", "", fmt.Errorf("%w: %w", ErrParseMACFailed, err)
	}

	if len(mac) != macBytes {
		return "", "", fmt.Errorf("%w: got %d bytes", ErrInvalidMACLength, len(mac))
	}

	interfaceID := generateInterfaceID(mac)
	if prefixStr == "" {
		return interfaceID, "", nil
	}

	prefixParts, err := parsePrefix(prefixStr)
	if err != nil {
		return "", "", err
	}

	fullIP := constructFullIP(prefixParts, mac)

	return interfaceID, fullIP, nil
}

// generateInterfaceID creates the EUI-64 interface ID from a MAC address.
func generateInterfaceID(mac []byte) string {
	eui64 := make([]byte, eui64Bytes)
	copy(eui64[0:3], mac[0:3])
	eui64[3] = fffeMarkerLow
	eui64[4] = fffeMarkerHigh
	copy(eui64[5:], mac[3:])
	eui64[0] ^= 0x02

	return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		eui64[0], eui64[1], eui64[2], eui64[3],
		eui64[4], eui64[5], eui64[6], eui64[7])
}

// parsePrefix validates and splits the IPv6 prefix.
func parsePrefix(prefixStr string) ([]string, error) {
	prefixStr = strings.TrimSuffix(prefixStr, "::")
	prefixParts := strings.Split(prefixStr, ":")

	if len(prefixParts) > PrefixMaxHextets {
		return nil, fmt.Errorf("%w: got %d hextets", ErrPrefixTooLong, len(prefixParts))
	}

	for i, part := range prefixParts {
		if part == "" {
			if i != 0 && i != len(prefixParts)-1 {
				return nil, fmt.Errorf("%w: %s", ErrInvalidEmptyHextet, prefixStr)
			}

			continue
		}

		var temp uint16
		if _, err := fmt.Sscanf(part, "%x", &temp); err != nil {
			return nil, fmt.Errorf("%w: %q: %w", ErrInvalidHextet, part, err)
		}
	}

	return prefixParts, nil
}

// constructFullIP builds the full IPv6 address from prefix parts and MAC.
func constructFullIP(prefixParts []string, mac []byte) string {
	ip6 := make([]uint16, ipv6Hextets)
	for hextetID := range ip6 {
		switch {
		case hextetID < len(prefixParts) && prefixParts[hextetID] != "":
			_, _ = fmt.Sscanf(prefixParts[hextetID], "%x", &ip6[hextetID]) // Safe—validated in parsePrefix
		case hextetID < PrefixMaxHextets:
			ip6[hextetID] = 0
		default:
			ip6[hextetID] = fillEUI64Hextet(hextetID, mac)
		}
	}

	return IP6ToString(ip6)
}

// fillEUI64Hextet computes EUI-64 hextet values from MAC bytes.
func fillEUI64Hextet(hextetID int, mac []byte) uint16 {
	eui64 := make([]byte, eui64Bytes)
	copy(eui64[0:3], mac[0:3])
	eui64[3] = fffeMarkerLow
	eui64[4] = fffeMarkerHigh
	copy(eui64[5:], mac[3:])
	eui64[0] ^= 0x02

	switch hextetID {
	case hextetIndex4:
		return uint16(eui64[0])<<byteShift | uint16(eui64[1])
	case hextetIndex5:
		return uint16(eui64[2])<<byteShift | uint16(eui64[3])
	case hextetIndex6:
		return uint16(eui64[4])<<byteShift | uint16(eui64[5])
	case hextetIndex7:
		return uint16(eui64[6])<<byteShift | uint16(eui64[7])
	}

	return 0 // Unreachable
}
