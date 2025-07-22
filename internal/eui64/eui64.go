// Package eui64 provides functionality for calculating EUI-64 interface identifiers and full IPv6 addresses from MAC addresses and prefixes.
// It includes the Calculator interface and a default implementation using the standard EUI-64 algorithm,
// along with helper functions for parsing, conversion, and string formatting of IPv6 addresses.
package eui64

import (
	"errors"
	"fmt"
	"net"
	"strconv"
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

	// IPv6 hextet indices for EUI-64 insertion.
	hextetEUI64First  = 4 // First hextet index where EUI-64 bytes are inserted.
	hextetEUI64Second = 5 // Second hextet index for EUI-64 bytes.
	hextetEUI64Third  = 6 // Third hextet index for EUI-64 bytes.
	hextetEUI64Fourth = 7 // Fourth hextet index for EUI-64 bytes.

	// Bit shift constant for combining bytes into a 16-bit hextet.
	byteShift = 8 // Number of bits to shift a byte to form a uint16.
)

// Static error variables.
var (
	ErrParseMAC             = errors.New("parsing MAC address")
	ErrInvalidMACLength     = fmt.Errorf("MAC address must be %d bytes", macBytes)
	ErrPrefixExceedsHextets = fmt.Errorf("IPv6 prefix exceeds %d hextets", prefixMaxHextets)
	ErrInvalidEmptyHextet   = errors.New("invalid empty hextet in IPv6 prefix")
)

// Calculator defines the interface for computing EUI-64 identifiers and IPv6 addresses.
type Calculator interface {
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
		return "", "", fmt.Errorf("%w: %w", ErrParseMAC, err)
	}

	if len(mac) != macBytes {
		return "", "", fmt.Errorf("%w, got %d", ErrInvalidMACLength, len(mac))
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
		return "", "", fmt.Errorf("%w, got %d", ErrPrefixExceedsHextets, len(prefixParts))
	}

	for i, part := range prefixParts {
		if part == "" && i != 0 && i != len(prefixParts)-1 {
			return "", "", fmt.Errorf("%w: %s", ErrInvalidEmptyHextet, prefixStr)
		}
	}

	ip6 := make([]uint16, ipv6Hextets)
	for i := range ip6 {
		switch {
		case i < len(prefixParts) && prefixParts[i] != "":
			if _, err := fmt.Sscanf(prefixParts[i], "%x", &ip6[i]); err != nil {
				return "", "", fmt.Errorf(
					"invalid hextet %q in IPv6 prefix: %w",
					prefixParts[i],
					err,
				)
			}
		case i < prefixMaxHextets:
			ip6[i] = 0
		case i == hextetEUI64First:
			ip6[i] = uint16(eui64[0])<<byteShift | uint16(eui64[1])
		case i == hextetEUI64Second:
			ip6[i] = uint16(eui64[2])<<byteShift | uint16(eui64[3])
		case i == hextetEUI64Third:
			ip6[i] = uint16(eui64[4])<<byteShift | uint16(eui64[5])
		case i == hextetEUI64Fourth:
			ip6[i] = uint16(eui64[6])<<byteShift | uint16(eui64[7])
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

	var stringBuilder strings.Builder

	prevWasCompression := false

	for i := 0; i < len(ip6); i++ {
		if i == bestStart && bestLen > 1 {
			stringBuilder.WriteString("::")

			prevWasCompression = true
			i += bestLen - 1

			continue
		}

		if ip6[i] != 0 || (i < bestStart || i >= bestStart+bestLen) {
			if i > 0 && !prevWasCompression {
				stringBuilder.WriteByte(':')
			}

			stringBuilder.WriteString(strconv.FormatUint(uint64(ip6[i]), 16))

			prevWasCompression = false
		}
	}

	return stringBuilder.String()
}
