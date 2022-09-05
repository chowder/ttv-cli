package redeemcustomreward

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSuccessResponse(t *testing.T) {
	bytes := []byte(`{
	  "data": {
		"redeemCommunityPointsCustomReward": {
		  "error": null,
		  "__typename": "RedeemCommunityPointsCustomRewardPayload"
		}
	  },
	  "extensions": {
		"durationMilliseconds": 201,
		"operationName": "RedeemCustomReward",
		"requestID": "<requestID>"
	  }
	}`)

	var response Response
	err := json.Unmarshal(bytes, &response)

	require.NoError(t, err)
	assert.Nil(t, response.Data.RedeemCommunityPointsCustomReward.Error)
}

func TestRewardOnCooldownResponse(t *testing.T) {
	bytes := []byte(`{
	  "data": {
		"redeemCommunityPointsCustomReward": {
		  "error": {
			"code": "GLOBAL_COOLDOWN",
			"__typename": "RedeemCommunityPointsCustomRewardError"
		  },
		  "__typename": "RedeemCommunityPointsCustomRewardPayload"
		}
	  },
	  "extensions": {
		"durationMilliseconds": 59,
		"operationName": "RedeemCustomReward",
		"requestID": "<requestID>"
	  }
	}`)

	var response Response
	err := json.Unmarshal(bytes, &response)

	require.NoError(t, err)
	assert.Equal(t, "GLOBAL_COOLDOWN", response.Data.RedeemCommunityPointsCustomReward.Error.Code)
}

func TestNotEnoughPointsResponse(t *testing.T) {
	bytes := []byte(`{
	  "data": {
		"redeemCommunityPointsCustomReward": {
		  "error": {
			"code": "NOT_ENOUGH_POINTS",
			"__typename": "RedeemCommunityPointsCustomRewardError"
		  },
		  "__typename": "RedeemCommunityPointsCustomRewardPayload"
		}
	  },
	  "extensions": {
		"durationMilliseconds": 35,
		"operationName": "RedeemCustomReward",
		"requestID": "<requestID>"
	  }
	}`)

	var response Response
	err := json.Unmarshal(bytes, &response)

	require.NoError(t, err)
	assert.Equal(t, "NOT_ENOUGH_POINTS", response.Data.RedeemCommunityPointsCustomReward.Error.Code)
}
