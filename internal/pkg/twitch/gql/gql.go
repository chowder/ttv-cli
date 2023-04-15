package gql

import (
	"encoding/json"
	"fmt"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch"
)

func Post(config config.Config, request any) ([]byte, error) {
	client, err := config.GetRestyClient()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request struct: %w", err)
	}

	resp, err := client.R().
		EnableTrace().
		SetHeader("Content-Type", "text/plain;charset=UTF-8").
		SetBody(body).
		Post(twitch.GqlApiUrl)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		// TODO: Narrow this to just the error within the body
		return nil, fmt.Errorf("HTTP request returned status code: %d - body: %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}
