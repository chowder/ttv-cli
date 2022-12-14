package gql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch"
)

func post(config *config.Config, request any) ([]byte, error) {
	// Make request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling GQL request: %w", err)
	}

	// Make a POST request
	req, err := http.NewRequest(http.MethodPost, twitch.GqlApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Client-ID", twitch.DefaultClientId)
	req.Header.Set("User-Agent", twitch.DefaultUserAgent)
	req.Header.Set("X-App-Version", twitch.DefaultAppVersion)
	req.Header.Set("Accept", twitch.DefaultAcceptHeader)
	req.Header.Set("Api-Consumer-Type", twitch.DefaultApiConsumerType)
	if config != nil {
		req.Header.Set("Client-Session-Id", config.GetClientSessionId())
		req.Header.Set("Authorization", "OAuth "+config.GetAuthToken())
		req.Header.Set("X-Device-Id", config.GetDeviceId())
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
	return post(nil, request)
}

func PostWithAuth(config *config.Config, request any) ([]byte, error) {
	return post(config, request)
}
