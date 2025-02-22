package tests

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
)

func TestValidateIPv6Prefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prefix  string
		wantErr error
	}{
		{
			name:    "Empty prefix",
			prefix:  "",
			wantErr: validators.ErrPrefixRequired,
		},
		{
			name:    "Blank prefix",
			prefix:  "   ",
			wantErr: validators.ErrPrefixRequired,
		},
		{
			name:    "Too many hextets",
			prefix:  "2001:0Db8:0000:0000:1234",
			wantErr: validators.ErrPrefixTooLong,
		},
		{
			name:    "Invalid empty hextet",
			prefix:  "2001::1234",
			wantErr: validators.ErrEmptyHextet,
		},
		{
			name:    "Invalid character",
			prefix:  "2001:0db8:gggg",
			wantErr: validators.ErrInvalidHextetChar,
		},
		{
			name:    "Too long hextet",
			prefix:  "2001:0db81234",
			wantErr: validators.ErrInvalidHextetLength,
		},
		{
			name:    "Valid prefix",
			prefix:  "2001:0db8:1234",
			wantErr: nil,
		},
		{
			name:    "Valid prefix with compression",
			prefix:  "2001:0db8::",
			wantErr: nil,
		},
	}

	validator := &validators.CombinedValidator{}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := validator.ValidateIPv6Prefix(testCase.prefix)
			assertValidationError(t, err, testCase.wantErr)
		})
	}
}
