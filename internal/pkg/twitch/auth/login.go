package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
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
			return "", fmt.Errorf("could not read username from terminal: %w", err)
		}
	}

	if len(password) == 0 {
		fmt.Print("Twitch password: ")
		//goland:noinspection GoRedundantConversion - This type cast is required for Windows
		b, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return "", fmt.Errorf("could not read password from terminal: %w", err)
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
		return "", fmt.Errorf("could not marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", loginApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("could not create HTTP request: %w", err)
	}
	req.Header.Add("Client-ID", twitch.DefaultClientId)

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing HTTP request: %w", err)
	}

	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading HTTP response body: %w", err)
	}

	if httpResp.StatusCode == 400 {
		type badRequestResponse struct {
			CaptchaProof string `json:"captcha_proof"`
			ErrorCode    int    `json:"error_code"`
		}

		var resp badRequestResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return "", fmt.Errorf("could not unmarshall request body: %w (body: %s)", err, string(body))
		}

		// Authy 2FA required
		if resp.ErrorCode == 3011 || resp.ErrorCode == 3012 {
			return loginWithAuthy2FA(resp.CaptchaProof, r)
		}

		if resp.ErrorCode == 1000 {
			return "", errors.New("captcha required for login - this is currently not supported")
		}

		return "", errors.New("Failed to login: " + string(body))
	}

	var result result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w (body: %s)", err, string(body))
	}

	return result.AccessToken, nil
}

func loginWithAuthy2FA(captchaProof string, r request) (string, error) {
	r.Captcha.Proof = captchaProof

	fmt.Print("2FA token: ")
	if _, err := fmt.Scanln(&r.AuthyToken); err != nil {
		return "", errors.New("could not read 2FA token from terminal")
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("could not marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", loginApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing HTTP request: %w", err)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading HTTP response body: %w", err)
	}

	var result result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshalling HTTP response body: %w", err)
	}

	return result.AccessToken, nil
}
