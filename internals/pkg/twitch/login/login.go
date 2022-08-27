package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"ttv-cli/internals/pkg/twitch"
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

func GetAccessToken(username string, password string) string {
	if len(username) == 0 {
		fmt.Print("Twitch username: ")
		if _, err := fmt.Scanln(&username); err != nil {
			log.Fatalln(err)
		}
	}

	if len(password) == 0 {
		fmt.Print("Twitch password: ")
		if _, err := fmt.Scanln(&password); err != nil {
			log.Fatalln(err)
		}
	}

	r := request{
		Username:     username,
		Password:     password,
		ClientId:     twitch.DefaultClientId,
		UndeleteUser: false,
		RememberMe:   true,
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", loginApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Client-ID", twitch.DefaultClientId)

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if httpResp.StatusCode == 400 {
		type badRequestResponse struct {
			CaptchaProof string `json:"captcha_proof"`
			ErrorCode    int    `json:"error_code"`
		}

		var resp badRequestResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			log.Fatalln(err)
		}

		// Authy 2FA required
		if resp.ErrorCode == 3011 || resp.ErrorCode == 3012 {
			return loginWithAuthy2FA(resp.CaptchaProof, r)
		}

		if resp.ErrorCode == 1000 {
			log.Fatalln("Captcha required for login - this is currently not supported")
		}

		log.Fatalln("Failed to login: ", string(body))
	}

	var result result
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}

	return result.AccessToken
}

func loginWithAuthy2FA(captchaProof string, r request) string {
	r.Captcha.Proof = captchaProof

	fmt.Print("2FA token: ")
	if _, err := fmt.Scanln(&r.AuthyToken); err != nil {
		log.Fatalln(err)
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", loginApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result result
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}
	return result.AccessToken
}
