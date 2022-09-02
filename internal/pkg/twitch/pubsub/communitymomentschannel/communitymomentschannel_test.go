package communitymomentschannel

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommunityPointsUser(t *testing.T) {
	bytes := []byte(`{
	  "type": "MESSAGE",
	  "data": {
		"topic": "community-moments-channel-v1.012345678",
		"message": "{\"type\":\"active\",\"data\":{\"moment_id\":\"7e8edd74-6db7-4299-b70a-5daf1f712e05\"}}"
	  }
	}`)

	var response Response
	err := json.Unmarshal(bytes, &response)

	require.NoError(t, err, string(bytes))

	assert.Equal(t, "MESSAGE", response.Type)
	assert.Equal(t, "community-moments-channel-v1.012345678", response.Data.Topic)

	require.NotEmpty(t, response.Data.Message)

	var message Message
	err = json.Unmarshal([]byte(response.Data.Message), &message)

	require.NoError(t, err, response.Data.Message)

	assert.Equal(t, "active", message.Type)
	require.NotEmpty(t, message.Data)

	var data ActiveMomentData
	err = json.Unmarshal(message.Data, &data)

	assert.Equal(t, "7e8edd74-6db7-4299-b70a-5daf1f712e05", data.MomentId)
}
