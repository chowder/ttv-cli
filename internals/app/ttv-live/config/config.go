package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Config struct {
	Streamers []string `json:"streamers"`
}

func GetDefaultConfigFile() string {
	home, _ := os.UserConfigDir()
	return path.Join(home, "ttv-live", "ttv-live.config")
}

func createDefaultConfig(filename string) Config {
	emptyConfig := Config{Streamers: []string{}}

	contents, err := json.MarshalIndent(emptyConfig, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(path.Base(filename), 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, contents, 0644); err != nil {
		log.Fatal(err)
	}

	return emptyConfig
}

func CreateOrReadFromFile(filename string) Config {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig(filename)
		}
		panic(err)
	}

	var config Config
	if err = json.Unmarshal(contents, &config); err != nil {
		panic(err)
	}
	return config
}
