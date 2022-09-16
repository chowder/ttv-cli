package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"time"
	"ttv-cli/internal/pkg/twitch"
	"ttv-cli/internal/pkg/utils"
)

var clientVersionPattern = regexp.MustCompile("window\\.__twilightBuildID=\"(.*?)\"")

type integrity struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

func (c *Config) GetIntegrityToken() (string, error) {
	err := c.refreshIntegrityToken()
	if err != nil {
		return "", fmt.Errorf("unable to refresh integrity token: %w", err)
	}

	return c.integrity.Token, nil
}

func (c *Config) refreshIntegrityToken() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Only refresh the token when it is less than 5 minutes from expiry
	if c.integrity != nil && time.Now().Add(time.Minute*5).Before(time.UnixMilli(c.integrity.Expiration)) {
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, twitch.IntegrityUrl, nil)
	if err != nil {
		return err
	}

	clientVersion, err := getClientVersion()
	if err != nil {
		return fmt.Errorf("error getting Twitch client version: %w", err)
	}

	clientSessionId, err := utils.TokenHex(16)
	if err != nil {
		return fmt.Errorf("error generating client session ID: %w", err)
	}

	req.Header.Set("Authorization", "OAuth "+c.authToken)
	req.Header.Set("Client-Id", twitch.DefaultClientId)
	req.Header.Set("Client-Session-Id", clientSessionId)
	req.Header.Set("User-Agent", twitch.DefaultUserAgent)
	req.Header.Set("Client-Version", clientVersion)
	req.Header.Set("X-Device-Id", c.deviceId)

	Config := &http.Client{}
	resp, err := Config.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &c.integrity)
	if err != nil {
		return err
	}

	c.clientVersion = clientVersion
	c.clientSessionId = clientSessionId

	return nil
}

func getClientVersion() (string, error) {
	resp, err := utils.HttpGet(twitch.TwitchHomeUrl, nil)
	if err != nil {
		return "", err
	}

	matches := clientVersionPattern.FindSubmatch(resp)
	if len(matches) < 2 {
		return "", errors.New("could not find client version")
	}

	return string(matches[1]), nil
}

func createRandomDeviceId() string {
	const chars = "abcdefghijklmnopqrstuvwxyz01234567890"
	const length = 26

	deviceId := make([]byte, length)
	for i := range deviceId {
		deviceId[i] = chars[rand.Intn(len(chars))]
	}

	return string(deviceId)
}
