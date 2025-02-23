package eui64

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/stretchr/testify/assert"
)

func TestConstructFullIP(t *testing.T) {
	tests := []struct {
		name        string
		prefixParts []string
		mac         []byte
		wantIP      string
	}{
		{
			name:        "Full prefix",
			prefixParts: []string{"2001", "0db8", "85a3", "0000"},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "2001:db8:85a3:0:214:22ff:fe01:2345",
		},
		{
			name:        "Partial prefix with empty",
			prefixParts: []string{"2001", "", "85a3"},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "2001:0:85a3:0:214:22ff:fe01:2345",
		},
		{
			name:        "Single hextet",
			prefixParts: []string{"2001"},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "2001::214:22ff:fe01:2345",
		},
		{
			name:        "Empty prefix parts",
			prefixParts: []string{""},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "::214:22ff:fe01:2345",
		},
		{
			name:        "Three hextets with trailing empty",
			prefixParts: []string{"2001", "db8", ""},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "2001:db8::214:22ff:fe01:2345",
		},
		{
			name:        "Prefix with mixed empty hextets",
			prefixParts: []string{"2001", "", "0db8", ""},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "2001:0:db8:0:214:22ff:fe01:2345", // Covers 115-121, 132-148
		},
		{
			name:        "Malformed hextet boundary",
			prefixParts: []string{"2001", "abcd"},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "2001:abcd::214:22ff:fe01:2345", // Covers 124-126
		},
		{
			name:        "Prefix with leading empty",
			prefixParts: []string{"", "2001", "0db8"},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "0:2001:db8:0:214:22ff:fe01:2345",
		},
		{
			name:        "Prefix with all empty",
			prefixParts: []string{"", "", "", ""},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "::214:22ff:fe01:2345",
		},
		{
			name:        "Prefix with single empty hextet",
			prefixParts: []string{"2001", "", "0db8", "abcd"},
			mac:         []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			wantIP:      "2001:0:db8:abcd:214:22ff:fe01:2345",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := eui64.ConstructFullIP(testCase.prefixParts, testCase.mac)
			assert.Equal(t, testCase.wantIP, got)
		})
	}
}

func TestFillEUI64Hextet(t *testing.T) {
	mac := []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}
	tests := []struct {
		name     string
		hextetID int
		want     uint16
	}{
		{name: "Hextet 4", hextetID: 4, want: 0x0214},         // 02 (flipped) << 8 | 14
		{name: "Hextet 5", hextetID: 5, want: 0x22FF},         // 22 << 8 | FF
		{name: "Hextet 6", hextetID: 6, want: 0xFE01},         // FE << 8 | 01
		{name: "Hextet 7", hextetID: 7, want: 0x2345},         // 23 << 8 | 45
		{name: "Invalid hextet ID 0", hextetID: 0, want: 0},   // Default case
		{name: "Invalid hextet ID -1", hextetID: -1, want: 0}, // Default case
		{name: "Invalid hextet ID 8", hextetID: 8, want: 0},   // Default case
		{name: "Invalid hextet ID 3", hextetID: 3, want: 0},   // Default case
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := eui64.FillEUI64Hextet(testCase.hextetID, mac)
			assert.Equal(t, testCase.want, got)
		})
	}
}
