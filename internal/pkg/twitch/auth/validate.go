package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const validateUrl = "https://id.twitch.tv/oauth2/validate"

type invalidTokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	ClientId  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserId    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
}

func Validate(token string) (Response, error) {
	if len(token) == 0 {
		return Response{}, errors.New("token cannot be null or empty")
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(validateUrl)

	if resp.StatusCode() == 401 {
		var response invalidTokenResponse
		if err = json.Unmarshal(resp.Body(), &response); err != nil {
			return Response{}, fmt.Errorf("error unmarshalling HTTP response body: %w", err)
		}

		return Response{}, errors.New(response.Message)
	}

	var response Response
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return response, fmt.Errorf("error unmarshalling HTTP response body: %w", err)
	}

	return response, nil
}
