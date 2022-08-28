package gql

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"ttv-cli/internal/pkg/twitch"
)

const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36"

func post(request any, authToken string) ([]byte, error) {
	// Make request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Make a POST request
	req, err := http.NewRequest("POST", twitch.GqlApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", twitch.DefaultClientId)
	req.Header.Set("User-Agent", defaultUserAgent)
	if len(authToken) > 0 {
		req.Header.Set("Authorization", "OAuth "+authToken)
	}

	// Execute the POST request
	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return body, nil
}

func Post(request any) ([]byte, error) {
	return post(request, "")
}

func PostWithAuth(request any, authToken string) ([]byte, error) {
	return post(request, authToken)
}
