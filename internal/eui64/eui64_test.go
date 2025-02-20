package eui64

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCalculateEUI64Valid tests the CalculateEUI64 function with valid inputs.
// It ensures that the function correctly computes the EUI-64 interface ID and full IPv6 address
// for various MAC addresses and IPv6 prefix combinations, including cases with no prefix and zero compression.
func TestCalculateEUI64Valid(t *testing.T) {
	tests := []struct {
		name            string
		mac             string
		prefix          string
		wantInterfaceID string
		wantFullIP      string
	}{
		{
			name:            "Valid MAC and full prefix",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:0db8:85a3:0000",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "2001:db8:85a3:0:214:22ff:fe01:2345",
		},
		{
			name:            "Valid MAC and partial prefix",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:0db8",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "2001:db8::214:22ff:fe01:2345",
		},
		{
			name:            "Valid MAC with no prefix",
			mac:             "00-14-22-01-23-45",
			prefix:          "",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "",
		},
		{
			name:            "Valid MAC with zero compression in prefix",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:0db8:0000:0000",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "2001:db8::214:22ff:fe01:2345",
		},
		{
			name:            "Valid MAC with trailing zero compression",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:db8::",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "2001:db8::214:22ff:fe01:2345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interfaceID, fullIP, err := CalculateEUI64(tt.mac, tt.prefix)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInterfaceID, interfaceID)
			assert.Equal(t, tt.wantFullIP, fullIP)
		})
	}
}

// TestCalculateEUI64Invalid tests the CalculateEUI64 function with invalid inputs.
// It verifies that the function returns appropriate errors for malformed MAC addresses and IPv6 prefixes,
// ensuring empty results when validation fails.
func TestCalculateEUI64Invalid(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		prefix  string
		wantErr string
	}{
		{
			name:    "Invalid MAC format",
			mac:     "invalid-mac",
			prefix:  "2001:0db8:85a3:0000",
			wantErr: "parsing MAC address",
		},
		{
			name:    "MAC too short",
			mac:     "00-14-22-01-23",
			prefix:  "2001:0db8:85a3:0000",
			wantErr: "parsing MAC address",
		},
		{
			name:    "Invalid prefix - too many hextets",
			mac:     "00-14-22-01-23-45",
			prefix:  "2001:0db8:85a3:0000:0000",
			wantErr: fmt.Sprintf("IPv6 prefix exceeds %d hextets", prefixMaxHextets),
		},
		{
			name:    "Invalid prefix - empty hextet",
			mac:     "00-14-22-01-23-45",
			prefix:  "2001::85a3",
			wantErr: "invalid empty hextet",
		},
		{
			name:    "Invalid prefix - invalid hextet",
			mac:     "00-14-22-01-23-45",
			prefix:  "2001:invalid:85a3",
			wantErr: "invalid hextet \"invalid\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interfaceID, fullIP, err := CalculateEUI64(tt.mac, tt.prefix)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.Empty(t, interfaceID)
			assert.Empty(t, fullIP)
		})
	}
}

// TestIp6ToString tests the ip6ToString function with various IPv6 address inputs.
// It ensures correct string formatting, including zero compression, for full addresses,
// addresses with zero runs, and edge cases like all zeros or leading/trailing zeros.
func TestIp6ToString(t *testing.T) {
	tests := []struct {
		name string
		ip6  []uint16
		want string
	}{
		{
			name: "No compression",
			ip6:  []uint16{0x2001, 0x0db8, 0x85a3, 0x0001, 0x0214, 0x22ff, 0xfe01, 0x2345},
			want: "2001:db8:85a3:1:214:22ff:fe01:2345",
		},
		{
			name: "Single zero compression",
			ip6:  []uint16{0x2001, 0x0db8, 0x0000, 0x0000, 0x0214, 0x22ff, 0xfe01, 0x2345},
			want: "2001:db8::214:22ff:fe01:2345",
		},
		{
			name: "All zeros",
			ip6:  []uint16{0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000},
			want: "::",
		},
		{
			name: "Leading zeros",
			ip6:  []uint16{0x0000, 0x0000, 0x0000, 0x0000, 0x0214, 0x22ff, 0xfe01, 0x2345},
			want: "::214:22ff:fe01:2345",
		},
		{
			name: "Trailing zeros",
			ip6:  []uint16{0x2001, 0x0db8, 0x85a3, 0x0001, 0x0000, 0x0000, 0x0000, 0x0000},
			want: "2001:db8:85a3:1::",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ip6ToString(tt.ip6))
		})
	}
}
