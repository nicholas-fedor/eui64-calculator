package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateMAC tests the ValidateMAC function with various MAC address inputs.
// It ensures the function accepts valid 48-bit MAC addresses and rejects malformed or incorrect ones,
// verifying error messages for invalid cases.
func TestValidateMAC(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantErr string
	}{
		{"Valid MAC", "00-14-22-01-23-45", ""},
		{"Empty MAC", "", "MAC address is required"},
		{"Invalid MAC format", "invalid-mac", "parsing MAC address"},
		{"MAC too short", "00-14-22-01-23", "parsing MAC address"},
		{"MAC with too many parts", "00-14-22-01-23-45-67", "parsing MAC address"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMAC(tt.mac)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
