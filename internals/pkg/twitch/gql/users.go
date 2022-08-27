package gql

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"ttv-cli/internals/pkg/twitch"
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

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a POST request
	req, err := http.NewRequest("POST", twitch.GqlApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Client-ID", twitch.DefaultClientId)

	// Execute the POST request
	client := &http.Client{}
	httpResp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(httpResp.Body)

	// Read the response body
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalln(err)
	}

	return response.Data.Users
}
