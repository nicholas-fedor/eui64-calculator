package eui64

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/stretchr/testify/require"
)

func TestCalculateEUI64(t *testing.T) {
	tests := []struct {
		name            string
		mac             string
		prefix          string
		wantInterfaceID string
		wantFullIP      string
		wantErr         string
	}{
		// Valid Cases
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
			name:            "Valid MAC with zero compression",
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
		{
			name:            "Valid MAC with single hextet",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "2001::214:22ff:fe01:2345",
		},
		{
			name:            "Valid MAC with empty hextet start",
			mac:             "00-14-22-01-23-45",
			prefix:          ":0db8:85a3",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "0:db8:85a3:0:214:22ff:fe01:2345",
		},
		{
			name:            "Valid MAC with three hextets",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:db8:85a3",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "2001:db8:85a3:0:214:22ff:fe01:2345",
		},
		// Invalid Cases
		{
			name:            "Invalid MAC format",
			mac:             "invalid-mac",
			prefix:          "2001:0db8:85a3:0000",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "parsing MAC address failed",
		},
		{
			name:            "MAC too short",
			mac:             "00-14-22-01-23",
			prefix:          "2001:0db8:85a3:0000",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "parsing MAC address failed",
		},
		{
			name:            "MAC too long",
			mac:             "00-14-22-01-23-45-67",
			prefix:          "2001:0db8:85a3:0000",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "parsing MAC address failed",
		},
		{
			name:            "Invalid prefix - too many hextets",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:0db8:85a3:0000:0000",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "IPv6 prefix exceeds 4 hextets",
		},
		{
			name:            "Invalid prefix - empty hextet middle",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001::85a3",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "invalid empty hextet",
		},
		{
			name:            "Invalid prefix - invalid hextet",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:invalid:85a3",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "invalid hextet in IPv6 prefix: \"invalid\"",
		},
		{
			name:            "Invalid prefix - bad hex char",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:0db8:xyz",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "invalid hextet in IPv6 prefix: \"xyz\"",
		},
		{
			name:            "MAC invalid length",
			mac:             "00-14-22-01-23-FF-FF",
			prefix:          "2001:0db8",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "parsing MAC address failed",
		},
		{
			name:            "Prefix with invalid middle empty",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001::85a3",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "invalid empty hextet",
		},
		{
			name:            "Invalid MAC parse",
			mac:             "invalid-mac",
			prefix:          "2001:0db8",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "parsing MAC address failed",
		},
		{
			name:            "GenerateInterfaceID only",
			mac:             "00-14-22-01-23-45",
			prefix:          "",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "",
		},
		{
			name:            "Invalid MAC length",
			mac:             "00-14-22-01-23-45-67",
			prefix:          "2001:0db8",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "parsing MAC address failed",
		},
		{
			name:            "Invalid prefix too long",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:0db8:85a3:0000:1234",
			wantInterfaceID: "",
			wantFullIP:      "",
			wantErr:         "IPv6 prefix exceeds 4 hextets",
		},
		{
			name:            "Prefix with full IP",
			mac:             "00-14-22-01-23-45",
			prefix:          "2001:0db8:85a3:0000",
			wantInterfaceID: "0214:22ff:fe01:2345",
			wantFullIP:      "2001:db8:85a3:0:214:22ff:fe01:2345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			calc := &eui64.DefaultCalculator{}
			interfaceID, fullIP, err := calc.CalculateEUI64(tt.mac, tt.prefix)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Empty(t, interfaceID)
				require.Empty(t, fullIP)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantInterfaceID, interfaceID)
				require.Equal(t, tt.wantFullIP, fullIP)
			}
		})
	}
}

func TestGenerateInterfaceID(t *testing.T) {
	tests := []struct {
		name string
		mac  []byte
		want string
	}{
		{
			name: "Standard MAC",
			mac:  []byte{0x00, 0x14, 0x22, 0x01, 0x23, 0x45},
			want: "0214:22ff:fe01:2345",
		},
		{
			name: "Different MAC",
			mac:  []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
			want: "a8bb:ccff:fedd:eeff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := eui64.GenerateInterfaceID(tt.mac)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParsePrefix(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		wantParts []string
		wantErr   string
	}{
		{
			name:      "Valid full prefix",
			prefix:    "2001:0db8:85a3:0000",
			wantParts: []string{"2001", "0db8", "85a3", "0000"},
		},
		{
			name:      "Valid partial prefix",
			prefix:    "2001:0db8",
			wantParts: []string{"2001", "0db8"},
		},
		{
			name:      "Valid with empty start",
			prefix:    ":0db8:85a3",
			wantParts: []string{"", "0db8", "85a3"},
		},
		{
			name:      "Valid with empty end",
			prefix:    "2001:0db8:",
			wantParts: []string{"2001", "0db8", ""},
		},
		{
			name:      "Too many hextets",
			prefix:    "2001:0db8:85a3:0000:0000",
			wantParts: nil,
			wantErr:   "IPv6 prefix exceeds 4 hextets",
		},
		{
			name:      "Empty hextet in middle",
			prefix:    "2001::85a3",
			wantParts: nil,
			wantErr:   "invalid empty hextet",
		},
		{
			name:      "Invalid hextet",
			prefix:    "2001:invalid:85a3",
			wantParts: nil,
			wantErr:   "invalid hextet in IPv6 prefix: \"invalid\"",
		},
		{
			name:      "Invalid hextet with letters",
			prefix:    "2001:xyz:85a3",
			wantParts: nil,
			wantErr:   "invalid hextet in IPv6 prefix: \"xyz\"",
		},
		{
			name:      "Empty hextet in middle",
			prefix:    "2001::85a3",
			wantParts: nil,
			wantErr:   "invalid empty hextet", // Covers 105-108
		},
		{
			name:      "Invalid hextet char",
			prefix:    "2001:0db8:xyz",
			wantParts: nil,
			wantErr:   "invalid hextet in IPv6 prefix", // Covers 96-100, 102
		},
		{
			name:      "Prefix with trailing colon",
			prefix:    "2001:0db8::", // Already tested, but ensure distinct
			wantParts: []string{"2001", "0db8"},
			wantErr:   "", // Covers 88-94 if not already
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			parts, err := eui64.ParsePrefix(testCase.prefix)
			if testCase.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), testCase.wantErr)
				require.Nil(t, parts)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.wantParts, parts)
			}
		})
	}
}
