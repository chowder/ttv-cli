package channelpointscontext

import (
	"encoding/json"
	"fmt"
	"ttv-cli/internal/pkg/twitch"
	"ttv-cli/internal/pkg/twitch/gql"
)

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type variables struct {
	ChannelLogin string `json:"channelLogin"`
}

type request struct {
	OperationName string     `json:"operationName"`
	Variables     variables  `json:"variables"`
	Extensions    extensions `json:"extensions"`
}

func makeRequest(channelLogin string) request {
	return request{
		OperationName: "ChannelPointsContext",
		Variables: variables{
			ChannelLogin: channelLogin,
		},
		Extensions: extensions{
			PersistedQuery: persistedQuery{
				Version:    1,
				Sha256Hash: "9988086babc615a918a1e9a722ff41d98847acac822645209ac7379eecb27152",
			},
		},
	}
}

func Get(client *twitch.Client, channelLogin string) (Response, error) {
	req := makeRequest(channelLogin)
	resp, err := gql.PostWithAuth(client, req)
	if err != nil {
		return Response{}, fmt.Errorf("error with GQL request: %w", err)
	}

	var c Response
	if err := json.Unmarshal(resp, &c); err != nil {
		return c, fmt.Errorf("error unmarshalling GQL response: %w", err)
	}

	return c, err
}
