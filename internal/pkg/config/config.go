package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"ttv-cli/internal/pkg/twitch/auth"
)

type TokenDetails struct {
	ClientId  string
	Login     string
	Scopes    []string
	UserId    string
	ExpiresIn int
}

type Config struct {
	authToken    string
	tokenDetails *TokenDetails
	deviceId     string
	integrity    *integrity
	mutex        sync.Mutex
}

type configJson struct {
	AuthToken string `json:"auth_token"`
	DeviceId  string `json:"device_id"`
}

func FromToken(authToken string) *Config {
	return &Config{
		authToken:    authToken,
		tokenDetails: nil,
		deviceId:     createRandomDeviceId(),
		integrity:    nil,
		mutex:        sync.Mutex{},
	}
}

func GetConfigFilePath() string {
	home, _ := os.UserConfigDir()
	return path.Join(home, "ttv-cli", "ttv-cli.config")
}

func CreateOrRead() (*Config, error) {
	configFilePath := GetConfigFilePath()
	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			authToken, err := auth.GetAccessToken("", "")
			if err != nil {
				return nil, fmt.Errorf("error getting Twitch auth token: %w", err)
			}
			c := FromToken(authToken)
			err = c.Save()
			if err != nil {
				return nil, fmt.Errorf("error saving Twitch auth token: %w", err)
			}
			return c, nil
		}
		return nil, fmt.Errorf("CreateOrRead: error when reading config file: %w", err)
	}

	var configJson configJson
	if err = json.Unmarshal(contents, &configJson); err != nil {
		return nil, fmt.Errorf("error when unmarshalling config file: %w", err)
	}

	c := configJson.ToConfig()
	if err := c.validateAuthToken(); err != nil {
		return nil, fmt.Errorf("error when validating Twitch auth token: %w", err)
	}

	return c, nil
}

func (c configJson) ToConfig() *Config {
	if len(c.DeviceId) == 0 {
		c.DeviceId = createRandomDeviceId()
	}

	return &Config{
		authToken:    c.AuthToken,
		tokenDetails: nil,
		deviceId:     c.DeviceId,
		integrity:    nil,
		mutex:        sync.Mutex{},
	}
}

func (c *Config) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	configFilePath := GetConfigFilePath()
	contents, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(configFilePath), 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	if err := os.WriteFile(configFilePath, contents, 0644); err != nil {
		return fmt.Errorf("could not write to config file: %w", err)
	}

	return nil
}

func (c *Config) GetAuthToken() string {
	return c.authToken
}

func (c *Config) GetDeviceId() string {
	return c.deviceId
}
