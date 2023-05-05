package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

var configFolderName = "ttv-cli"
var configFileName = "config.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	return path.Join(home, configFolderName, configFileName), nil
}

func (config Config) Save() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("could not get config file path: %w", err)
	}

	err = os.MkdirAll(path.Dir(configFilePath), 0755)
	if err != nil {
		return fmt.Errorf("could not create config path: %w", err)
	}

	contents, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	err = os.WriteFile(configFilePath, contents, 0644)
	if err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}

func Load() (Config, error) {
	var config Config

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("could not get config file path: %w", err)
	}

	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, fmt.Errorf("could not read from config file: %w", err)
	}

	err = json.Unmarshal(contents, &config)
	if err != nil {
		return config, fmt.Errorf("could not unmarshall config file contents: %w", err)
	}

	return config, nil
}
