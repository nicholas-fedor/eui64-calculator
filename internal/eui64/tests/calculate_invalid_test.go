package tests

import (
	"fmt"
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
)

func TestCalculateEUI64Invalid(t *testing.T) {
	t.Parallel()

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
			wantErr: fmt.Sprintf("IPv6 prefix exceeds %d hextets", eui64.PrefixMaxHextets),
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
			wantErr: "invalid hextet in IPv6 prefix: \"invalid\"",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			interfaceID, fullIP, err := eui64.CalculateEUI64(testCase.mac, testCase.prefix)
			assertEUI64Result(t, interfaceID, fullIP, "", "", err, testCase.wantErr)
		})
	}
}
