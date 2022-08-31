package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const validateUrl = "https://id.twitch.tv/oauth2/validate"

type invalidTokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Validate(token string) error {
	client := &http.Client{}

	httpReq, err := http.NewRequest("GET", validateUrl, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "OAuth "+token)

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %w", err)
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode == 401 {
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return fmt.Errorf("error reading HTTP response body: %w", err)
		}

		var resp invalidTokenResponse
		if err = json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("error unmarshalling HTTP response body: %w", err)
		}

		return errors.New(resp.Message)
	}

	return nil
}
