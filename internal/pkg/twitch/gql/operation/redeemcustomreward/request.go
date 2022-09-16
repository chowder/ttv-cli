package redeemcustomreward

import (
	"encoding/json"
	"fmt"
	"ttv-cli/internal/pkg/config"
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

func makeRequest(input Input) (request, error) {
	token, err := utils.TokenHex(16)
	if err != nil {
		return request{}, fmt.Errorf("could not generate transaction ID: %w", err)
	}

	input.TransactionID = token
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
	}, nil
}

func Redeem(c *config.Config, input Input) (Response, error) {
	var response Response
	req, err := makeRequest(input)
	if err != nil {
		return response, fmt.Errorf("error generating GQL request: %w", err)
	}

	body, err := gql.PostWithAuth(c, req)
	if err != nil {
		return response, fmt.Errorf("error with GQL request: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("could not unmarshal response: %s, error: %w", string(body), err)
	}

	return response, nil
}
