package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateIPv6Prefix tests the ValidateIPv6Prefix function with various prefix inputs.
// It verifies that the function accepts valid IPv6 prefixes (up to 4 hextets) with proper formatting
// and rejects invalid ones, checking for specific error messages on failure.
func TestValidateIPv6Prefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		wantErr string
	}{
		{"Valid full IPv6 prefix", "2001:db8:85a3:0", ""},
		{"Valid partial IPv6 prefix", "2001:db8", ""},
		{"Valid prefix with trailing ::", "2001:db8::", ""},
		{"Valid minimal prefix", "::", ""},
		{"Invalid blank prefix", "", "a non-blank IPv6 prefix is expected"},
		{"Invalid character in hextet", "2001:db8:85a3:g000", "invalid character in hextet"},
		{"IPv4 address", "192.168.1.1", "invalid character in hextet"},
		{"Too many hextets", "2001:db8:85a3:0:0", fmt.Sprintf("IPv6 prefix must be %d or fewer hextets", maxHextets)},
		{"Invalid internal empty hextet", "2001::85a3:0", "empty hextet in IPv6 prefix"},
		{"Invalid hextet length", "2001:db8:85a3:12345", "invalid hextet length in IPv6 prefix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPv6Prefix(tt.prefix)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
