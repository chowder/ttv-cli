package login

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateSuccess(t *testing.T) {
	bytes := []byte(`{
	  "client_id": "dummyclientid",
	  "login": "forsen",
	  "scopes": [
		"chat_login",
		"user_presence_friends_read",
		"user_read"
	  ],
	  "user_id": "012345678",
	  "expires_in": 0
	}`)

	var response Response
	err := json.Unmarshal(bytes, &response)

	require.NoError(t, err)

	assert.Equal(t, "dummyclientid", response.ClientId)
	assert.Equal(t, "forsen", response.Login)
	assert.ElementsMatch(t, []string{"chat_login", "user_presence_friends_read", "user_read"}, response.Scopes)
	assert.Equal(t, "012345678", response.UserId)
	assert.Equal(t, 0, response.ExpiresIn)
}
