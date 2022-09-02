package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	client := &http.Client{}

	httpReq, err := http.NewRequest("GET", validateUrl, nil)
	if err != nil {
		return Response{}, fmt.Errorf("error creating HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "OAuth "+token)

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("error performing HTTP request: %w", err)
	}

	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("error reading HTTP response body: %w", err)
	}

	if httpResp.StatusCode == 401 {
		var response invalidTokenResponse
		if err = json.Unmarshal(body, &response); err != nil {
			return Response{}, fmt.Errorf("error unmarshalling HTTP response body: %w", err)
		}

		return Response{}, errors.New(response.Message)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("error unmarshalling HTTP response body: %w", err)
	}

	return response, nil
}
