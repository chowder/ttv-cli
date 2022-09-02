package communitypointsuser

// topic: community-points-user-v1

import (
	"encoding/json"
	"time"
)

// PointsSpentData Use when Response#Type is 'points-spent'
type PointsSpentData struct {
	Timestamp time.Time `json:"timestamp"`
	Balance   struct {
		UserId    string `json:"user_id"`
		ChannelId string `json:"channel_id"`
		Balance   int    `json:"balance"`
	} `json:"balance"`
}

// PointsEarnedData Use when Response#Type is 'points-earned'
type PointsEarnedData struct {
	Timestamp time.Time `json:"timestamp"`
	ChannelId string    `json:"channel_id"`
	PointGain struct {
		UserId         string `json:"user_id"`
		ChannelId      string `json:"channel_id"`
		TotalPoints    int    `json:"total_points"`
		BaselinePoints int    `json:"baseline_points"`
		ReasonCode     string `json:"reason_code"`
		Multipliers    []struct {
			ReasonCode string  `json:"reason_code"`
			Factor     float64 `json:"factor"`
		} `json:"multipliers"`
	} `json:"point_gain"`
	Balance struct {
		UserId    string `json:"user_id"`
		ChannelId string `json:"channel_id"`
		Balance   int    `json:"balance"`
	} `json:"balance"`
}

type Response struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
