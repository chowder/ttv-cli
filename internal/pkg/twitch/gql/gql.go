package gql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ttv-cli/internal/pkg/twitch"
)

const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36"

func post(request any, authToken string) ([]byte, error) {
	// Make request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling GQL request: %w", err)
	}

	// Make a POST request
	req, err := http.NewRequest("POST", twitch.GqlApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
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
		return nil, fmt.Errorf("error performing HTTP request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(httpResp.Body)

	// Read the response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP response body: %w", err)
	}

	if httpResp.StatusCode != 200 {
		// TODO: Narrow this to just the error within the body
		return nil, fmt.Errorf("HTTP request returned status code: %d - body: %s", httpResp.StatusCode, string(body))
	}

	return body, nil
}

func Post(request any) ([]byte, error) {
	return post(request, "")
}

func PostWithAuth(request any, authToken string) ([]byte, error) {
	return post(request, authToken)
}
