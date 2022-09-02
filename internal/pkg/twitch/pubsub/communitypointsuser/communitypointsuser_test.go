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
	  "type": "MESSAGE",
	  "data": {
		"topic": "community-points-user-v1.012345678",
		"message": "{\"type\":\"points-spent\",\"data\":{\"timestamp\":\"2022-02-02T22:22:22.020202020Z\",\"balance\":{\"user_id\":\"012345678\",\"channel_id\":\"987654321\",\"balance\":3820}}}"
	  }
	}`)

	var response Response
	err := json.Unmarshal(pointsSpentResponse, &response)

	require.NoError(t, err, string(pointsSpentResponse))

	assert.Equal(t, "MESSAGE", response.Type)
	assert.Equal(t, "community-points-user-v1.012345678", response.Data.Topic)

	require.NotEmpty(t, response.Data.Message)

	var message Message
	err = json.Unmarshal([]byte(response.Data.Message), &message)

	require.NoError(t, err, response.Data.Message)

	assert.Equal(t, "points-spent", message.Type)
	assert.NotEmpty(t, message.Data)

	var pointsEarned PointsSpentData
	err = json.Unmarshal(message.Data, &pointsEarned)

	require.NoError(t, err)

	expectedTimestamp := time.Date(2022, 2, 2, 22, 22, 22, 20202020, time.UTC)
	assert.Equal(t, expectedTimestamp, pointsEarned.Timestamp)
	assert.Equal(t, "012345678", pointsEarned.Balance.UserId)
	assert.Equal(t, "987654321", pointsEarned.Balance.ChannelId)
	assert.Equal(t, 3820, pointsEarned.Balance.Balance)
}

func TestPointsEarnedResponse(t *testing.T) {
	pointsEarnedResponse := []byte(`{
	  "type": "MESSAGE",
	  "data": {
		"topic": "community-points-user-v1.012345678",
		"message": "{\"type\":\"points-earned\",\"data\":{\"timestamp\":\"2022-09-01T22:12:23.698533302Z\",\"channel_id\":\"987654321\",\"point_gain\":{\"user_id\":\"012345678\",\"channel_id\":\"987654321\",\"total_points\":12,\"baseline_points\":10,\"reason_code\":\"WATCH\",\"multipliers\":[{\"reason_code\":\"SUB_T1\",\"factor\":0.2}]},\"balance\":{\"user_id\":\"012345678\",\"channel_id\":\"987654321\",\"balance\":624285}}}"
	  }
	}`)

	var response Response
	err := json.Unmarshal(pointsEarnedResponse, &response)

	require.NoError(t, err, string(pointsEarnedResponse))

	assert.Equal(t, "MESSAGE", response.Type)
	assert.Equal(t, "community-points-user-v1.012345678", response.Data.Topic)

	require.NotEmpty(t, response.Data.Message)

	var message Message
	err = json.Unmarshal([]byte(response.Data.Message), &message)

	require.NoError(t, err, response.Data.Message)

	assert.Equal(t, "points-earned", message.Type)
	assert.NotEmpty(t, message.Data)

	var pointsEarned PointsEarnedData
	err = json.Unmarshal(message.Data, &pointsEarned)

	require.NoError(t, err, string(message.Data))

	// Point Gain
	pointGain := pointsEarned.PointGain
	assert.Equal(t, "012345678", pointGain.UserId)
	assert.Equal(t, "987654321", pointGain.ChannelId)
	assert.Equal(t, 12, pointGain.TotalPoints)
	assert.Equal(t, 10, pointGain.BaselinePoints)

	multiplier := pointsEarned.PointGain.Multipliers[0]
	assert.Equal(t, "SUB_T1", multiplier.ReasonCode)
	assert.Equal(t, 0.2, multiplier.Factor)

	// Balance
	balance := pointsEarned.Balance
	assert.Equal(t, "012345678", balance.UserId)
	assert.Equal(t, "987654321", balance.ChannelId)
	assert.Equal(t, 624285, balance.Balance)
}
