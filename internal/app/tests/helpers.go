package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// assertRunError checks the error result of the Run function.
func assertRunError(t *testing.T, err error, wantErr bool, expectedErr error) {
	t.Helper()

	if wantErr {
		require.Error(t, err)

		if expectedErr != nil {
			require.Equal(t, expectedErr, err)
		}
	} else {
		require.NoError(t, err)
	}
}
