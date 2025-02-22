package tests

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/stretchr/testify/assert"
)

func TestConstructFullIP(t *testing.T) {
	t.Parallel()

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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := eui64.ConstructFullIP(tc.prefixParts, tc.mac)
			assert.Equal(t, tc.wantIP, got)
		})
	}
}

func TestFillEUI64Hextet(t *testing.T) {
	t.Parallel()

	mac := []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}
	tests := []struct {
		name     string
		hextetID int
		want     uint16
	}{
		{name: "Hextet 4", hextetID: 4, want: 0x0214}, // 02 (flipped) << 8 | 14
		{name: "Hextet 5", hextetID: 5, want: 0x22FF}, // 22 << 8 | FF
		{name: "Hextet 6", hextetID: 6, want: 0xFE01}, // FE << 8 | 01
		{name: "Hextet 7", hextetID: 7, want: 0x2345}, // 23 << 8 | 45
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := eui64.FillEUI64Hextet(tc.hextetID, mac)
			assert.Equal(t, tc.want, got)
		})
	}
}
