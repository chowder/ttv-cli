package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/term"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall"
	"ttv-cli/internal/pkg/twitch"
)

const loginApiUrl = "https://passport.twitch.tv/login"

type request struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	UndeleteUser bool   `json:"undelete_user"`
	RememberMe   bool   `json:"remember_me"`
	Captcha      struct {
		Proof string `json:"proof"`
	} `json:"captcha,omitempty"`
	AuthyToken string `json:"authy_token,omitempty"`
}

type result struct {
	AccessToken string `json:"access_token"`
}

func GetAccessToken(username string, password string) (string, error) {
	if len(username) == 0 {
		fmt.Print("Twitch username: ")
		if _, err := fmt.Scanln(&username); err != nil {
			return "", err
		}
	}

	if len(password) == 0 {
		fmt.Print("Twitch password: ")
		b, err := term.ReadPassword(syscall.Stdin)
		fmt.Println()
		if err != nil {
			return "", err
		}
		password = string(b)
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	r := request{
		Username:     username,
		Password:     password,
		ClientId:     twitch.DefaultClientId,
		UndeleteUser: false,
		RememberMe:   true,
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", loginApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Add("Client-ID", twitch.DefaultClientId)

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	if httpResp.StatusCode == 400 {
		type badRequestResponse struct {
			CaptchaProof string `json:"captcha_proof"`
			ErrorCode    int    `json:"error_code"`
		}

		var resp badRequestResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return "", err
		}

		// Authy 2FA required
		if resp.ErrorCode == 3011 || resp.ErrorCode == 3012 {
			return loginWithAuthy2FA(resp.CaptchaProof, r)
		}

		if resp.ErrorCode == 1000 {
			return "", errors.New("Captcha required for login - this is currently not supported")
		}

		return "", errors.New(fmt.Sprint("Failed to login: ", string(body)))
	}

	var result result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func loginWithAuthy2FA(captchaProof string, r request) (string, error) {
	r.Captcha.Proof = captchaProof

	fmt.Print("2FA token: ")
	if _, err := fmt.Scanln(&r.AuthyToken); err != nil {
		return "", err
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", loginApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	var result result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
