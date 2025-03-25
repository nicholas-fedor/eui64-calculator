package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateMAC tests the ValidateMAC function with various MAC address inputs.
// It ensures the function accepts valid 48-bit MAC addresses and rejects empty, overly long, or malformed inputs,
// verifying error messages for each validation step.
func TestValidateMAC(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantErr string
	}{
		// Valid cases
		{"Valid MAC with hyphens", "00-14-22-01-23-45", ""},
		{"Valid MAC with colons", "00:14:22:01:23:45", ""},

		// Empty check
		{"Empty MAC", "", "MAC address is required"},
		{"Whitespace-only MAC", "   ", "MAC address is required"},

		// Length check
		{
			"MAC just over max length",
			"00-14-22-01-23-45-6",
			fmt.Sprintf("MAC address string exceeds maximum length of %d characters", macStrLen),
		},
		{
			"MAC exceeds max length",
			"00-14-22-01-23-45-6789",
			fmt.Sprintf("MAC address string exceeds maximum length of %d characters", macStrLen),
		},
		{
			"MAC with too many parts",
			"00-14-22-01-23-45-67",
			fmt.Sprintf("MAC address string exceeds maximum length of %d characters", macStrLen),
		},

		// Parsing errors
		{"Invalid MAC format (non-hex)", "invalid-mac", "parsing MAC address"},
		{"MAC too short", "00-14-22-01-23", "parsing MAC address"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMAC(tt.mac)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
