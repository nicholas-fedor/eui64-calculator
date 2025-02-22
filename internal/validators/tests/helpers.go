package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// assertValidationError checks if the validation error contains the expected error message.
func assertValidationError(t *testing.T, err error, wantErr error) {
	t.Helper()

	if wantErr != nil {
		require.Error(t, err)
		require.Contains(t, err.Error(), wantErr.Error(), "Error message mismatch")
	} else {
		require.NoError(t, err)
	}
}
