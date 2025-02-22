package mocks

import (
	"errors"
	"testing"
)

func TestValidator_ValidateMAC(t *testing.T) {
	tests := []struct {
		name    string
		macErr  error
		wantErr string // Changed to string for error message comparison
	}{
		{"No error", nil, ""},
		{"Invalid MAC", ErrInvalidMAC, "invalid MAC"},
		{"Generic error", errors.New("test error"), "test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{MacErr: tt.macErr}
			err := v.ValidateMAC("00:11:22:33:44:55")

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != tt.wantErr {
				t.Errorf("ValidateMAC() error = %v, wantErr %v", errMsg, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateIPv6Prefix(t *testing.T) {
	tests := []struct {
		name      string
		prefixErr error
		wantErr   string // Changed to string for error message comparison
	}{
		{"No error", nil, ""},
		{"Invalid prefix", ErrInvalidPrefix, "invalid prefix"},
		{"Generic error", errors.New("test error"), "test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{PrefixErr: tt.prefixErr}
			err := v.ValidateIPv6Prefix("2001:db8::/32")

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != tt.wantErr {
				t.Errorf("ValidateIPv6Prefix() error = %v, wantErr %v", errMsg, tt.wantErr)
			}
		})
	}
}
