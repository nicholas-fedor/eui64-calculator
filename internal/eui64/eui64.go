package eui64

import (
	"fmt"
	"net"
	"strings"
)

// Constants defining sizes and markers for EUI-64 and IPv6 calculations.
const (
	ipv6Hextets      = 8    // ipv6Hextets is the number of hextets in a full IPv6 address.
	prefixMaxHextets = 4    // prefixMaxHextets is the maximum number of hextets allowed in an IPv6 prefix.
	macBytes         = 6    // macBytes is the expected length of a MAC address in bytes.
	eui64Bytes       = 8    // eui64Bytes is the length of an EUI-64 identifier in bytes.
	fffeMarkerLow    = 0xFF // fffeMarkerLow is the low byte of the EUI-64 FFFE marker.
	fffeMarkerHigh   = 0xFE // fffeMarkerHigh is the high byte of the EUI-64 FFFE marker.
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

// CalculateEUI64 computes the EUI-64 interface ID and full IPv6 address from a MAC address and an optional IPv6 prefix.
// It converts the MAC address to an EUI-64 identifier by inserting the FFFE marker and flipping the local/global bit,
// then constructs an IPv6 address if a prefix is provided. Returns the interface ID, full IPv6 address, and any error.
func CalculateEUI64(macStr, prefixStr string) (string, string, error) {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return "", "", fmt.Errorf("parsing MAC address: %w", err)
	}
	if len(mac) != macBytes {
		return "", "", fmt.Errorf("MAC address must be %d bytes, got %d", macBytes, len(mac))
	}

	eui64 := make([]byte, eui64Bytes)
	copy(eui64[0:3], mac[0:3])
	eui64[3] = fffeMarkerLow
	eui64[4] = fffeMarkerHigh
	copy(eui64[5:], mac[3:])
	eui64[0] ^= 0x02 // Flip the local/global bit (7th bit) for EUI-64.

	interfaceID := fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		eui64[0], eui64[1], eui64[2], eui64[3],
		eui64[4], eui64[5], eui64[6], eui64[7])

	if prefixStr == "" {
		return interfaceID, "", nil
	}

	prefixStr = strings.TrimSuffix(prefixStr, "::")
	prefixParts := strings.Split(prefixStr, ":")
	if len(prefixParts) > prefixMaxHextets {
		return "", "", fmt.Errorf("IPv6 prefix exceeds %d hextets, got %d", prefixMaxHextets, len(prefixParts))
	}

	for i, part := range prefixParts {
		if part == "" && i != 0 && i != len(prefixParts)-1 {
			return "", "", fmt.Errorf("invalid empty hextet in IPv6 prefix: %s", prefixStr)
		}
	}

	ip6 := make([]uint16, ipv6Hextets)
	for i := range ip6 {
		if i < len(prefixParts) && prefixParts[i] != "" {
			if _, err := fmt.Sscanf(prefixParts[i], "%x", &ip6[i]); err != nil {
				return "", "", fmt.Errorf("invalid hextet %q in IPv6 prefix: %w", prefixParts[i], err)
			}
		} else if i < prefixMaxHextets {
			ip6[i] = 0
		} else {
			switch i {
			case 4:
				ip6[i] = uint16(eui64[0])<<8 | uint16(eui64[1])
			case 5:
				ip6[i] = uint16(eui64[2])<<8 | uint16(eui64[3])
			case 6:
				ip6[i] = uint16(eui64[4])<<8 | uint16(eui64[5])
			case 7:
				ip6[i] = uint16(eui64[6])<<8 | uint16(eui64[7])
			}
		}
	}

	return interfaceID, ip6ToString(ip6), nil
}

// ip6ToString converts a 128-bit IPv6 address into its canonical string representation.
// It applies zero compression (e.g., "::") to the longest run of consecutive zero hextets,
// ensuring a compact and valid IPv6 address format.
func ip6ToString(ip6 []uint16) string {
	allZeros := true
	for _, h := range ip6 {
		if h != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		return "::"
	}

	// Find the longest run of zeros for compression.
	bestStart, bestLen := -1, 0
	start, length := -1, 0
	for i, h := range ip6 {
		if h == 0 {
			if start == -1 {
				start = i
			}
			length++
			if length > bestLen && length > 1 {
				bestStart, bestLen = start, length
			}
		} else if start != -1 {
			start, length = -1, 0
		}
	}
	if start != -1 && length > bestLen && length > 1 {
		bestStart, bestLen = start, length
	}

	var b strings.Builder
	prevWasCompression := false
	for i := 0; i < len(ip6); i++ {
		if i == bestStart && bestLen > 1 {
			b.WriteString("::")
			prevWasCompression = true
			i += bestLen - 1
			continue
		}
		if ip6[i] != 0 || (i < bestStart || i >= bestStart+bestLen) {
			if i > 0 && !prevWasCompression {
				b.WriteByte(':')
			}
			b.WriteString(fmt.Sprintf("%x", ip6[i]))
			prevWasCompression = false
		}
	}
	return b.String()
}
