package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Config struct {
	AuthToken string `json:"auth_token"`
}

func GetConfigFilePath() string {
	home, _ := os.UserConfigDir()
	return path.Join(home, "ttv-cli", "ttv-cli.config")
}

func createDefaultConfig() Config {
	emptyConfig := Config{AuthToken: ""}
	emptyConfig.Save()
	return emptyConfig
}

func CreateOrRead() Config {
	configFilePath := GetConfigFilePath()
	contents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig()
		}
		log.Fatalln(err)
	}

	var config Config
	if err = json.Unmarshal(contents, &config); err != nil {
		log.Fatalln(err)
	}
	// TODO: Verify that config.AuthToken is valid
	return config
}

func (c Config) Save() {
	configFilePath := GetConfigFilePath()
	contents, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(path.Base(configFilePath), 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(configFilePath, contents, 0644); err != nil {
		log.Fatal(err)
	}
}
