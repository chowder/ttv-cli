package config

import (
	"errors"
	"fmt"
	"github.com/Adeithe/go-twitch/pubsub"
	"github.com/go-resty/resty/v2"
	"ttv-cli/internal/pkg/twitch"
	"ttv-cli/internal/pkg/twitch/auth"
	twitchUtils "ttv-cli/internal/pkg/twitch/utils"
	"ttv-cli/internal/pkg/utils"
)

type Config struct {
	Username        string `json:"username"`
	UserId          string `json:"user_id"`
	AuthToken       string `json:"auth_token"`
	DeviceId        string `json:"device_id"`
	ClientVersion   string `json:"client_version"`
	ClientSessionId string `json:"client_session_id"`
}

func Create() (Config, error) {
	config, err := CreateNoAuth()
	if err != nil {
		return config, err
	}

	authToken, err := auth.GetAccessToken()
	if err != nil {
		return config, fmt.Errorf("error creating auth token: %w", err)
	}
	config.AuthToken = authToken

	resp, err := auth.Validate(authToken)
	if err != nil {
		return config, fmt.Errorf("error validating auth token: %w", err)
	}
	config.UserId = resp.UserId

	config.Username = resp.Login

	return config, nil
}

func CreateNoAuth() (Config, error) {
	var config Config

	config.DeviceId = auth.CreateRandomDeviceId()

	clientVersion, err := twitchUtils.GetClientVersion()
	if err != nil {
		return config, fmt.Errorf("error getting current Twitch client version: %w", err)
	}

	config.ClientVersion = clientVersion
	config.ClientSessionId, _ = utils.TokenHex(16)

	return config, nil
}

func (config Config) GetRestyClient() (*resty.Client, error) {
	client := resty.New()

	if len(config.AuthToken) > 0 {
		client = client.SetHeader("Authorization", "OAuth "+config.AuthToken)
	}

	client = client.
		SetHeader("Client-Id", twitch.DefaultClientId)

	if len(config.ClientSessionId) == 0 {
		return nil, errors.New("client session ID not created")
	}
	client = client.SetHeader("Client-Session-Id", config.ClientSessionId)

	if len(config.ClientVersion) == 0 {
		return nil, errors.New("client version not fetched")
	}
	client = client.SetHeader("Client-Version", config.ClientVersion)

	client = client.SetHeader("User-Agent", twitch.DefaultUserAgent)

	if len(config.DeviceId) == 0 {
		return nil, errors.New("device ID not created")
	}
	client = client.SetHeader("X-Device-Id", config.DeviceId)

	return client, nil
}

func (config Config) PubSubListen(client *pubsub.Client, topic string, args ...interface{}) error {
	return client.ListenWithAuth(config.AuthToken, topic, args...)
}
