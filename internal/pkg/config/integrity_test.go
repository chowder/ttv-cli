package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetClientVersion(t *testing.T) {
	version, err := getClientVersion()

	require.NoError(t, err)
	require.NotEmpty(t, version)
}