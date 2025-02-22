package eui64

import (
	"strconv"
	"strings"
)

// IP6ToString converts a 128-bit IPv6 address into its canonical string representation.
func IP6ToString(ip6 []uint16) string {
	if isAllZeros(ip6) {
		return "::"
	}

	bestStart, bestLen := findLongestZeroRun(ip6)

	return formatIPv6(ip6, bestStart, bestLen)
}

// isAllZeros checks if all hextets are zero.
func isAllZeros(ip6 []uint16) bool {
	for _, hextet := range ip6 {
		if hextet != 0 {
			return false
		}
	}

	return true
}

// findLongestZeroRun identifies the longest run of zeros for compression.
func findLongestZeroRun(ip6 []uint16) (int, int) {
	bestStart, bestLen := -1, 0
	currentStart, currentLen := -1, 0

	for i, hextet := range ip6 {
		currentStart, currentLen = updateZeroRun(i, hextet, currentStart, currentLen)
		bestStart, bestLen = updateBestRun(currentStart, currentLen, bestStart, bestLen)
	}

	return finalizeBestRun(currentStart, currentLen, bestStart, bestLen)
}

// updateZeroRun adjusts the current zero run based on the hextet value.
func updateZeroRun(index int, hextet uint16, start, length int) (int, int) {
	if hextet == 0 {
		if start == -1 {
			start = index
		}

		length++
	} else if start != -1 {
		start, length = -1, 0
	}

	return start, length
}

// updateBestRun updates the best zero run if the current one is longer.
func updateBestRun(currentStart, currentLen, bestStart, bestLen int) (int, int) {
	if currentLen > bestLen && currentLen > 1 {
		return currentStart, currentLen
	}

	return bestStart, bestLen
}

// finalizeBestRun checks if the last run is the longest.
func finalizeBestRun(currentStart, currentLen, bestStart, bestLen int) (int, int) {
	if currentStart != -1 && currentLen > bestLen && currentLen > 1 {
		return currentStart, currentLen
	}

	return bestStart, bestLen
}

// formatIPv6 builds the string with zero compression applied.
func formatIPv6(ip6 []uint16, bestStart, bestLen int) string {
	var builder strings.Builder

	prevWasCompression := false

	for ip6Digit := 0; ip6Digit < len(ip6); ip6Digit++ {
		if ip6Digit == bestStart && bestLen > 1 {
			builder.WriteString("::")

			prevWasCompression = true
			ip6Digit += bestLen - 1

			continue
		}

		if ip6[ip6Digit] != 0 || (ip6Digit < bestStart || ip6Digit >= bestStart+bestLen) {
			if ip6Digit > 0 && !prevWasCompression {
				builder.WriteByte(':')
			}

			builder.WriteString(strconv.FormatUint(uint64(ip6[ip6Digit]), 16))

			prevWasCompression = false
		}
	}

	return builder.String()
}
