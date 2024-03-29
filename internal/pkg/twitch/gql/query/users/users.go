package users

import (
	"encoding/json"
	"fmt"
	"time"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql"
)

const getUsersQuery = `query Users($logins: [String!]) {
    users(logins: $logins) {
        displayName
		id
		login
        profileURL
        stream {
            game {
                displayName
            }
            createdAt
            viewersCount
        }
    }
}`

type variables struct {
	Logins []string `json:"logins"`
}

type request struct {
	Query     string    `json:"query"`
	Variables variables `json:"variables"`
}

type Stream struct {
	Game struct {
		DisplayName string `json:"displayName"`
	}
	CreatedAt    time.Time `json:"createdAt"`
	ViewersCount int       `json:"viewersCount"`
}

type User struct {
	DisplayName string  `json:"displayName"`
	Id          string  `json:"id"`
	Login       string  `json:"login"`
	ProfileURL  string  `json:"profileURL"`
	Stream      *Stream `json:"stream"`
}

type Response struct {
	Data struct {
		Users []User `json:"users"`
	} `json:"data"`
}

func makeRequest(logins []string) request {
	return request{
		Query: getUsersQuery,
		Variables: variables{
			Logins: logins,
		},
	}
}

func GetUsers(config config.Config, logins []string) ([]User, error) {
	request := makeRequest(logins)

	gqlResp, err := gql.Post(config, request)
	if err != nil {
		return nil, fmt.Errorf("GetUsers: error with GQL request: %w", err)
	}

	var response Response
	if err := json.Unmarshal(gqlResp, &response); err != nil {
		return nil, fmt.Errorf("GetUsers: error unmarshalling GQL response: %w", err)
	}

	return response.Data.Users, nil
}
