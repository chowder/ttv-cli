package communitypointsuser

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPointsSpendResponse(t *testing.T) {
	pointsSpentResponse := []byte(`{
	  "type": "points-spent",
	  "data": {
		"timestamp": "2022-02-02T22:22:22.020202020Z",
		"balance": {
		  "user_id": "012345678",
		  "channel_id": "987654321",
		  "balance": 3820
		}
	  }
	}`)

	var response Response
	err := json.Unmarshal(pointsSpentResponse, &response)

	require.NoError(t, err)
	require.Equal(t, "points-spent", response.Type)
	require.NotEmpty(t, response.Data)

	var data PointsSpentData
	err = json.Unmarshal(response.Data, &data)

	require.NoError(t, err, string(response.Data))
	expectedTimestamp := time.Date(2022, 2, 2, 22, 22, 22, 20202020, time.UTC)
	assert.Equal(t, expectedTimestamp, data.Timestamp)
	assert.Equal(t, "012345678", data.Balance.UserId)
	assert.Equal(t, "987654321", data.Balance.ChannelId)
	assert.Equal(t, 3820, data.Balance.Balance)
}

func TestPointsEarnedResponse(t *testing.T) {
	pointsEarnedResponse := []byte(`{
	  "type": "points-earned",
	  "data": {
		"timestamp": "2022-09-01T22:12:23.698533302Z",
		"channel_id": "987654321",
		"point_gain": {
		  "user_id": "012345678",
		  "channel_id": "987654321",
		  "total_points": 12,
		  "baseline_points": 10,
		  "reason_code": "WATCH",
		  "multipliers": [
			{
			  "reason_code": "SUB_T1",
			  "factor": 0.2
			}
		  ]
		},
		"balance": {
		  "user_id": "012345678",
		  "channel_id": "987654321",
		  "balance": 624285
		}
	  }
	}`)

	var response Response
	err := json.Unmarshal(pointsEarnedResponse, &response)

	require.NoError(t, err, string(pointsEarnedResponse))
	require.Equal(t, "points-earned", response.Type)
	require.NotEmpty(t, response.Data)

	var data PointsEarnedData
	err = json.Unmarshal(response.Data, &data)

	require.NoError(t, err, string(response.Data))

	// Point Gain
	pointGain := data.PointGain
	assert.Equal(t, "012345678", pointGain.UserId)
	assert.Equal(t, "987654321", pointGain.ChannelId)
	assert.Equal(t, 12, pointGain.TotalPoints)
	assert.Equal(t, 10, pointGain.BaselinePoints)

	multiplier := data.PointGain.Multipliers[0]
	assert.Equal(t, "SUB_T1", multiplier.ReasonCode)
	assert.Equal(t, 0.2, multiplier.Factor)

	// Balance
	balance := data.Balance
	assert.Equal(t, "012345678", balance.UserId)
	assert.Equal(t, "987654321", balance.ChannelId)
	assert.Equal(t, 624285, balance.Balance)
}
