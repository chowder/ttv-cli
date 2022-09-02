package communitymomentschannel

import "encoding/json"

type ActiveMomentData struct {
	MomentId string `json:"moment_id"`
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Response struct {
	Type string `json:"type"`
	Data struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	} `json:"data"`
}
