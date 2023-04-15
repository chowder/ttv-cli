package channelfollows

import (
	"encoding/json"
	"fmt"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql"
)

type variables struct {
	Limit int    `json:"limit"`
	Order string `json:"order"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type channelFollowsRequest struct {
	OperationName string     `json:"operationName"`
	Variables     variables  `json:"variables"`
	Extensions    extensions `json:"extensions"`
}

func makeRequest() channelFollowsRequest {
	return channelFollowsRequest{
		OperationName: "ChannelFollows",
		Variables: variables{
			Limit: 100,
			Order: "ASC",
		},
		Extensions: extensions{
			PersistedQuery: persistedQuery{
				Version:    1,
				Sha256Hash: "eecf815273d3d949e5cf0085cc5084cd8a1b5b7b6f7990cf43cb0beadf546907",
			},
		},
	}
}

type ChannelFollow struct {
	DisplayName string
	Login       string
}

// Get TODO: Implement cursor following to handle >100 follows
func Get(c config.Config) ([]ChannelFollow, error) {
	req := makeRequest()
	respBody, err := gql.Post(c, req)
	if err != nil {
		return nil, fmt.Errorf("error with GQL request: %w", err)
	}

	var resp response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling GQL request: %w", err)
	}

	follows := make([]ChannelFollow, 0)
	for _, node := range resp.Data.User.Follows.Edges {
		if node.Typename == "FollowEdge" {
			follows = append(follows, ChannelFollow{
				DisplayName: node.Node.DisplayName,
				Login:       node.Node.Login,
			})
		}
	}

	return follows, nil
}
