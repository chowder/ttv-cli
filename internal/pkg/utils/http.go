package utils

import (
	"fmt"
	"io"
	"net/http"
	"ttv-cli/internal/pkg/twitch"
)

func HttpGet(url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", twitch.DefaultUserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}
