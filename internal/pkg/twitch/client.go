package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Integrity struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

type Client struct {
	authToken string
	deviceId  string
	integrity Integrity
	mutex     sync.Mutex
}

func NewClient(authToken string) *Client {
	return &Client{
		authToken: authToken,
		deviceId:  createRandomDeviceId(),
		integrity: Integrity{},
		mutex:     sync.Mutex{},
	}
}

func (c Client) GetAuthToken() string {
	return c.authToken
}

func (c Client) GetDeviceId() string {
	return c.deviceId
}

func (c *Client) GetIntegrityToken() (string, error) {
	err := c.refreshIntegrityToken()
	if err != nil {
		return "", fmt.Errorf("unable to refresh integrity token: %w", err)
	}

	return c.integrity.Token, nil
}

func (c *Client) refreshIntegrityToken() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expiresAt := time.UnixMilli(c.integrity.Expiration)

	// Only refresh the token when it is less than 5 minutes from expiry
	if time.Now().Add(time.Minute * 5).Before(expiresAt) {
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, IntegrityUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "OAuth "+c.authToken)
	req.Header.Set("Client-Id", DefaultClientId)
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("X-Device-Id", c.deviceId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &c.integrity)
	if err != nil {
		return err
	}

	return nil
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
