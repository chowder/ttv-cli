package login

import (
	"encoding/json"
	"errors"
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
		return err
	}
	httpReq.Header.Set("Authorization", "OAuth "+token)

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode == 401 {
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		var resp invalidTokenResponse
		if err = json.Unmarshal(body, &resp); err != nil {
			return err
		}

		return errors.New(resp.Message)
	}

	return nil
}
