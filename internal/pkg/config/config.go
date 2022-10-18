package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"ttv-cli/internal/pkg/twitch/auth"
	"ttv-cli/internal/pkg/utils"
)

type TokenDetails struct {
	ClientId  string
	Login     string
	Scopes    []string
	UserId    string
	ExpiresIn int
}

type Config struct {
	authToken       string
	clientVersion   string
	clientSessionId string
	tokenDetails    *TokenDetails
	deviceId        string
	mutex           sync.Mutex
}

type configJson struct {
	AuthToken       string `json:"auth_token"`
	ClientSessionId string `json:"client_session_id"`
	DeviceId        string `json:"device_id"`
}

func FromToken(authToken string) *Config {
	clientSessionId, _ := utils.TokenHex(16)
	return &Config{
		authToken:       authToken,
		tokenDetails:    nil,
		deviceId:        createRandomDeviceId(),
		clientSessionId: clientSessionId,
		mutex:           sync.Mutex{},
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

	_ = c.Save()
	
	return c, nil
}

func (c configJson) ToConfig() *Config {
	if len(c.DeviceId) == 0 {
		c.DeviceId = createRandomDeviceId()
	}
	if len(c.ClientSessionId) == 0 {
		c.ClientSessionId, _ = utils.TokenHex(16)
	}

	return &Config{
		authToken:       c.AuthToken,
		tokenDetails:    nil,
		deviceId:        c.DeviceId,
		clientSessionId: c.ClientSessionId,
		mutex:           sync.Mutex{},
	}
}

func (c *Config) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	configFilePath := GetConfigFilePath()

	cj := configJson{
		AuthToken:       c.authToken,
		ClientSessionId: c.clientSessionId,
		DeviceId:        c.deviceId,
	}

	contents, err := json.MarshalIndent(cj, "", "  ")
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

func (c *Config) GetClientSessionId() string {
	return c.clientSessionId
}
