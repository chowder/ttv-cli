package communitymomentcalloutclaim

import (
	"fmt"
	"log"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql"
)

type input struct {
	MomentID string `json:"momentID"`
}

type variables struct {
	Input input `json:"input"`
}

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type request struct {
	OperationName string     `json:"operationName"`
	Variables     variables  `json:"variables"`
	Extensions    extensions `json:"extensions"`
}

func makeRequest(momentId string) request {
	return request{
		OperationName: "CommunityMomentCallout_Claim",
		Variables: variables{
			Input: input{
				MomentID: momentId,
			},
		},
		Extensions: extensions{
			PersistedQuery: persistedQuery{
				Version:    1,
				Sha256Hash: "e2d67415aead910f7f9ceb45a77b750a1e1d9622c936d832328a0689e054db62",
			},
		},
	}
}

func Claim(c *config.Config, momentId string) error {
	req := makeRequest(momentId)
	resp, err := gql.PostWithAuth(c, req)
	if err != nil {
		return fmt.Errorf("error with GQL request: %w", err)
	}

	log.Println(string(resp))
	return nil
}
