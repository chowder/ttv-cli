package spade

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetSpadeUrl(t *testing.T) {
	url, err := GetUrl("xqc")

	require.NoError(t, err)
	require.NotEmpty(t, url)
}
