package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"ttv-cli/internal/pkg/twitch/login"
)

type Config struct {
	AuthToken string `json:"auth_token"`
}

func GetConfigFilePath() string {
	home, _ := os.UserConfigDir()
	return path.Join(home, "ttv-cli", "ttv-cli.config")
}

func createDefaultConfig() (Config, error) {
	emptyConfig := Config{AuthToken: ""}

	err := emptyConfig.validateAuthToken()
	if err != nil {
		return Config{}, fmt.Errorf("createDefaultConfig: error when validating Twitch auth token: %w", err)
	}

	return emptyConfig, nil
}

func CreateOrRead() (Config, error) {
	configFilePath := GetConfigFilePath()
	contents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig()
		}
		return Config{}, fmt.Errorf("CreateOrRead: error when reading config file: %w", err)
	}

	var config Config
	if err = json.Unmarshal(contents, &config); err != nil {
		return Config{}, fmt.Errorf("CreateOrRead: Error when unmarshalling config file: %w", err)
	}

	if err = config.validateAuthToken(); err != nil {
		return Config{}, fmt.Errorf("CreateOrRead: Error when validating Twitch auth token: %w", err)
	}

	return config, nil
}

func (c Config) validateAuthToken() error {
	if len(c.AuthToken) == 0 || login.Validate(c.AuthToken) != nil {
		fmt.Println("Auth token not found or expired, generating a new one for you...")

		authToken, err := login.GetAccessToken("", "")
		if err != nil {
			return fmt.Errorf("validateAuthToken: Error getting Twitch access token: %w", err)
		}

		c.AuthToken = authToken
		if err := c.Save(); err != nil {
			return fmt.Errorf("validateAuthToken: Error saving config: %w", err)
		}
	}
	return nil
}

func (c Config) Save() error {
	configFilePath := GetConfigFilePath()
	contents, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(configFilePath), 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	if err := ioutil.WriteFile(configFilePath, contents, 0644); err != nil {
		return fmt.Errorf("could not write to config file: %w", err)
	}

	return nil
}
