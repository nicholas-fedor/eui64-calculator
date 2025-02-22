package validators

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
)

func TestValidateMAC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		mac     string
		name    string
		wantErr error
	}{
		{
			name:    "Empty MAC",
			mac:     "",
			wantErr: validators.ErrMACRequired,
		},
		{
			name:    "Empty MAC with spaces",
			mac:     "   ",
			wantErr: validators.ErrMACRequired,
		},
		{
			name:    "Invalid MAC format",
			mac:     "invalid-mac",
			wantErr: validators.ErrParseMACFailed,
		},
		{
			name:    "MAC too short",
			mac:     "00-14-22-01-23",
			wantErr: validators.ErrParseMACFailed,
		},
		{
			name:    "Valid MAC",
			mac:     "00-14-22-01-23-45",
			wantErr: nil,
		},
	}

	validator := &validators.CombinedValidator{}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := validator.ValidateMAC(testCase.mac)
			assertValidationError(t, err, testCase.wantErr)
		})
	}
}
