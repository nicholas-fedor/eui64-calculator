package tests

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIp6ToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ip6  []uint16
		want string
	}{
		{
			name: "All zeros",
			ip6:  []uint16{0, 0, 0, 0, 0, 0, 0, 0},
			want: "::",
		},
		{
			name: "No compression needed",
			ip6:  []uint16{0x2001, 0x0db8, 0x85a3, 0x0000, 0x0214, 0x22ff, 0xfe01, 0x2345},
			want: "2001:db8:85a3:0:214:22ff:fe01:2345",
		},
		{
			name: "Compression at start",
			ip6:  []uint16{0, 0, 0, 0, 0x0214, 0x22ff, 0xfe01, 0x2345},
			want: "::214:22ff:fe01:2345",
		},
		{
			name: "Compression in middle",
			ip6:  []uint16{0x2001, 0x0db8, 0, 0, 0, 0x22ff, 0xfe01, 0x2345},
			want: "2001:db8::22ff:fe01:2345",
		},
		{
			name: "Compression at end",
			ip6:  []uint16{0x2001, 0x0db8, 0x85a3, 0x0001, 0, 0, 0, 0},
			want: "2001:db8:85a3:1::",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := eui64.IP6ToString(testCase.ip6)
			require.NotNil(t, got)
			assert.Equal(t, testCase.want, got)
		})
	}
}
