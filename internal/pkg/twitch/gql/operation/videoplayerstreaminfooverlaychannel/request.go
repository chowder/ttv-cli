package videoplayerstreaminfooverlaychannel

import (
	"encoding/json"
	"fmt"
	"ttv-cli/internal/pkg/twitch/gql"
)

type persistedquery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type extensions struct {
	PersistedQuery persistedquery `json:"persistedQuery"`
}

type variables struct {
	Channel string `json:"channel"`
}

type request struct {
	OperationName string     `json:"operationName"`
	Extensions    extensions `json:"extensions"`
	Variables     variables  `json:"variables"`
}

func makeRequest(channel string) request {
	return request{
		OperationName: "VideoPlayerStreamInfoOverlayChannel",
		Extensions: extensions{
			PersistedQuery: persistedquery{
				Version:    1,
				Sha256Hash: "a5f2e34d626a9f4f5c0204f910bab2194948a9502089be558bb6e779a9e1b3d2",
			},
		},
		Variables: variables{
			Channel: channel,
		},
	}
}

func Get(channel string) (Response, error) {
	req := makeRequest(channel)
	body, err := gql.Post(req)
	if err != nil {
		return Response{}, fmt.Errorf("error with GQL request: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, fmt.Errorf("error unmarshalling GQL response: %s, error: %w", string(body), err)
	}

	return resp, nil
}
