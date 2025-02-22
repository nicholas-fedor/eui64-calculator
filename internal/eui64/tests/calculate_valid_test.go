package tests

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
)

func TestCalculateEUI64Valid(t *testing.T) {
	t.Parallel()

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
		{
			name:            "Valid MAC with single hextet prefix",
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
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			interfaceID, fullIP, err := eui64.CalculateEUI64(testCase.mac, testCase.prefix)
			assertEUI64Result(t, interfaceID, fullIP, testCase.wantInterfaceID, testCase.wantFullIP, err, "")
		})
	}
}
