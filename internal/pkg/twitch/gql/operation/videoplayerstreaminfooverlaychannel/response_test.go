package videoplayerstreaminfooverlaychannel

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResponse(t *testing.T) {
	resp, err := Get("xqc")

	require.NoError(t, err)
	require.NotEmpty(t, resp)
}
