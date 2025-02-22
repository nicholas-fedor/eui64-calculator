package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertEUI64Result checks the results of CalculateEUI64 for equality and error conditions.
func assertEUI64Result(
	t *testing.T,
	interfaceID,
	fullIP,
	wantInterfaceID,
	wantFullIP string,
	err error,
	wantErr string,
) {
	t.Helper()

	if wantErr != "" {
		require.Error(t, err)
		assert.Contains(t, err.Error(), wantErr)
		assert.Empty(t, interfaceID)
		assert.Empty(t, fullIP)
	} else {
		require.NoError(t, err)
		assert.Equal(t, wantInterfaceID, interfaceID)
		assert.Equal(t, wantFullIP, fullIP)
	}
}
