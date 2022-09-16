package channelfollows

import (
	"encoding/json"
	"fmt"
	"ttv-cli/internal/pkg/twitch"
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
				Sha256Hash: "4b9cb31b54b9213e5760f2f6e9e935ad09924cac2f78aac51f8a64d85f028ed0",
			},
		},
	}
}

type ChannelFollow struct {
	DisplayName string
	Login       string
}

// Get TODO: Implement cursor following to handle >100 follows
func Get(client *twitch.Client) ([]ChannelFollow, error) {
	req := makeRequest()
	respBody, err := gql.PostWithAuth(client, req)
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
