package users

import (
	"encoding/json"
	"log"
	"ttv-cli/internal/pkg/twitch/gql"
)

const getUsersQuery = `query Users($logins: [String!]) {
    users(logins: $logins) {
        displayName
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

type User struct {
	DisplayName string `json:"displayName"`
	ProfileURL  string `json:"profileURL"`
	Stream      struct {
		Game struct {
			DisplayName string `json:"displayName"`
		}
		CreatedAt    string `json:"createdAt"`
		ViewersCount int    `json:"viewersCount"`
	} `json:"stream"`
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

func GetUsers(logins []string) []User {
	// Create POST request body
	request := makeRequest(logins)

	gqlResp, err := gql.Post(request)
	if err != nil {
		log.Fatalln(err)
	}

	var response Response
	if err := json.Unmarshal(gqlResp, &response); err != nil {
		log.Fatalln(err)
	}

	return response.Data.Users
}
