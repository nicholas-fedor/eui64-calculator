package mocks

import (
	"errors"
	"testing"
)

func TestCalculator_CalculateEUI64(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		interfaceID string
		fullIP      string
		wantErr     string // Changed to string for error message comparison
	}{
		{"Success", nil, "interface", "full_ip", ""},
		{"With error", errors.New("calc error"), "", "", "calc error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Calculator{
				Err:         tt.err,
				InterfaceID: tt.interfaceID,
				FullIP:      tt.fullIP,
			}

			gotInterface, gotFullIP, err := m.CalculateEUI64("mac", "prefix")

			// Check error
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != tt.wantErr {
				t.Errorf("CalculateEUI64() error = %v, wantErr %v", errMsg, tt.wantErr)
			}

			if err == nil {
				if gotInterface != tt.interfaceID {
					t.Errorf("CalculateEUI64() interface = %v, want %v", gotInterface, tt.interfaceID)
				}

				if gotFullIP != tt.fullIP {
					t.Errorf("CalculateEUI64() fullIP = %v, want %v", gotFullIP, tt.fullIP)
				}
			}
		})
	}
}
