package redeemcustomreward

import (
	"ttv-cli/internal/pkg/twitch/gql"
	"ttv-cli/internal/pkg/utils"
)

type variables struct {
	Input Input `json:"input"`
}

type Input struct {
	ChannelID     string `json:"channelID"`
	Cost          int    `json:"cost"`
	Prompt        string `json:"prompt"`
	RewardID      string `json:"rewardID"`
	TextInput     string `json:"textInput,omitempty"`
	Title         string `json:"title"`
	TransactionID string `json:"transactionID"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type request struct {
	OperationName string     `json:"operationName"`
	Variables     variables  `json:"variables"`
	Extensions    extensions `json:"extensions"`
}

func makeRequest(input Input) request {
	input.TransactionID = utils.TokenHex(16)
	return request{
		OperationName: "RedeemCustomReward",
		Variables: variables{
			Input: input,
		},
		Extensions: extensions{
			PersistedQuery: persistedQuery{
				Version:    1,
				Sha256Hash: "d56249a7adb4978898ea3412e196688d4ac3cea1c0c2dfd65561d229ea5dcc42",
			},
		},
	}
}

func Redeem(input Input, authToken string) ([]byte, error) {
	req := makeRequest(input)
	return gql.PostWithAuth(req, authToken)
}
