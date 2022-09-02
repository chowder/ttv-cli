package communitymomentschannel

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommunityPointsUser(t *testing.T) {
	bytes := []byte(`{
	  "type": "active",
	  "data": {
		"moment_id": "7e8edd74-6db7-4299-b70a-5daf1f712e05"
	  }
	}`)

	var response Response
	err := json.Unmarshal(bytes, &response)

	require.NoError(t, err, string(bytes))

	assert.Equal(t, "active", response.Type)
	assert.Equal(t, "7e8edd74-6db7-4299-b70a-5daf1f712e05", response.Data.MomentId)
}
