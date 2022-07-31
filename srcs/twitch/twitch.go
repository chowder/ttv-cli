package twitch

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const QUERY = `query Channel($logins: [String!]) {
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

const GqlApiUri = "https://gql.twitch.tv/gql"

const DefaultClientId = "kimne78kx3ncx6brgo4mv6wki5h1ko"

type Variables struct {
	Logins []string `json:"logins"`
}

type Request struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables"`
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

func makeRequest(logins []string) Request {
	return Request{
		Query: QUERY,
		Variables: Variables{
			Logins: logins,
		},
	}
}

func GetAllStreamers(logins []string) []User {
	// Create POST request body
	request := makeRequest(logins)

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a POST request
	req, err := http.NewRequest("POST", GqlApiUri, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Client-ID", DefaultClientId)

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
