package eui64

import (
	"fmt"
	"net"
	"strings"
)

// CalculateEUI64 converts a MAC address and IPv6 prefix into an EUI-64 IPv6 address.
func CalculateEUI64(macStr, prefixStr string) (string, string, error) {
	// Parse MAC address
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return "", "", fmt.Errorf("invalid MAC address: %v", err)
	}

	// Ensure MAC is 48 bits (6 bytes)
	if len(mac) != 6 {
		return "", "", fmt.Errorf("MAC address must be 48 bits")
	}

	// Split MAC into two parts and insert FFFE
	eui64 := make([]byte, 8)
	copy(eui64[0:3], mac[0:3])
	eui64[3] = 0xFF
	eui64[4] = 0xFE
	copy(eui64[5:8], mac[3:6])

	// Flip the 7th bit (Universal/Local bit)
	eui64[0] ^= 0x02

	// Format the interface ID (end of IPv6 address)
	interfaceID := fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		eui64[0], eui64[1], eui64[2], eui64[3],
		eui64[4], eui64[5], eui64[6], eui64[7])

	// Handle the prefix
	if prefixStr == "" {
		// If no prefix is provided, return just the interface ID
		return interfaceID, "", nil
	}

	// Parse and validate the prefix
	prefixParts := strings.Split(prefixStr, ":")
	if len(prefixParts) > 4 {
		return "", "", fmt.Errorf("IPv6 prefix must be 4 or fewer hextets, got %d", len(prefixParts))
	}

	// Validate each part of the prefix
	for _, part := range prefixParts {
		if part == "" {
			return "", "", fmt.Errorf("empty hextet in IPv6 prefix: %s", prefixStr)
		}
		if len(part) > 4 {
			return "", "", fmt.Errorf("invalid hextet length in IPv6 prefix: %s", part)
		}
		_, err := fmt.Sscanf(part, "%x", new(uint16))
		if err != nil {
			return "", "", fmt.Errorf("invalid hextet in IPv6 prefix: %s", part)
		}
	}

	// Construct the full IPv6 address array (8 hextets)
	ip6 := make([]uint16, 8)
	for i := 0; i < 8; i++ {
		if i < len(prefixParts) {
			fmt.Sscanf(prefixParts[i], "%x", &ip6[i])
		} else if i < 4 {
			ip6[i] = 0 // Fill remaining prefix hextets with zeros
		} else {
			// Fill the interface ID part
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

	// Convert the IPv6 array to a string with zero compression
	fullIP := ip6ToString(ip6)

	return interfaceID, fullIP, nil
}

// ip6ToString converts an IPv6 address array to a string with zero compression.
func ip6ToString(ip6 []uint16) string {
	// Handle the special case of all zeros
	isAllZeros := true
	for _, hextet := range ip6 {
		if hextet != 0 {
			isAllZeros = false
			break
		}
	}
	if isAllZeros {
		return "::"
	}

	// Find the longest stretch of zeros
	var zeroStart, zeroLength, maxZeroStart, maxZeroLength int
	inZeroRun := false

	for i := 0; i < 8; i++ {
		if ip6[i] == 0 {
			if !inZeroRun {
				zeroStart = i
				zeroLength = 1
				inZeroRun = true
			} else {
				zeroLength++
			}
		} else {
			if inZeroRun && zeroLength > maxZeroLength && zeroLength > 1 {
				maxZeroStart = zeroStart
				maxZeroLength = zeroLength
			}
			inZeroRun = false
		}
	}
	if inZeroRun && zeroLength > maxZeroLength && zeroLength > 1 {
		maxZeroStart = zeroStart
		maxZeroLength = zeroLength
	}

	// Build the string without compression
	var parts []string
	for i := 0; i < 8; i++ {
		parts = append(parts, fmt.Sprintf("%x", ip6[i]))
	}

	// If no compression is needed, join and return
	if maxZeroLength <= 1 {
		return strings.Join(parts, ":")
	}

	// Apply zero compression
	var compressedParts []string
	for i := 0; i < 8; i++ {
		if i == maxZeroStart && maxZeroLength > 1 {
			compressedParts = append(compressedParts, "") // Insert empty part for zero compression
			i += maxZeroLength - 1                        // Skip the zero run
		} else {
			compressedParts = append(compressedParts, parts[i])
		}
	}

	// Join parts with colons
	result := strings.Join(compressedParts, ":")

	// Handle leading and trailing zeros for zero compression
	if maxZeroStart == 0 {
		result = ":" + result // Leading zeros
	}
	if maxZeroStart+maxZeroLength == 8 {
		result += ":" // Trailing zeros
	}

	return result
}
