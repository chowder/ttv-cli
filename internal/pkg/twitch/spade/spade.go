package spade

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"ttv-cli/internal/pkg/twitch"
	"ttv-cli/internal/pkg/utils"
)

var (
	spadeUrlPattern          = regexp.MustCompile("\"spade_url\":\"(.*?)\"")
	twitchSettingsUrlPattern = regexp.MustCompile("(https://static.twitchcdn.net/config/settings.*?js)")
)

func toTwitchUrl(streamer string) string {
	return fmt.Sprintf("https://twitch.tv/%s", strings.ToLower(streamer))
}

// GetUrl TODO: Find out if Spade URLs are specific to each streamer
func GetUrl(streamer string) (string, error) {
	url := toTwitchUrl(streamer)
	resp, err := utils.HttpGet(url, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching streamer URL: %s, error: %w", url, err)
	}

	settingsUrl := string(twitchSettingsUrlPattern.Find(resp))

	resp, err = utils.HttpGet(settingsUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching settings url: %s, error %w", settingsUrl, err)
	}

	matches := spadeUrlPattern.FindSubmatch(resp)
	if len(matches) < 2 {
		return "", errors.New("could not find Spade URL")
	}

	return string(matches[1]), nil
}

type properties struct {
	ChannelId   string `json:"channel_id"`
	BroadcastId string `json:"broadcast_id"`
	Player      string `json:"player"`
	UserId      string `json:"user_id"`
}

type data struct {
	Event      string     `json:"event"`
	Properties properties `json:"properties"`
}

type payload struct {
	Data []byte `json:"data"`
}

func SendWatchMinute(spadeUrl string, channelId string, broadcastId string, userId string) error {
	d := data{
		Event: "minute-watched",
		Properties: properties{
			ChannelId:   channelId,
			BroadcastId: broadcastId,
			Player:      "site",
			UserId:      userId,
		},
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	p := payload{
		Data: b,
	}

	b, err = json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, spadeUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", twitch.DefaultUserAgent)

	c := &http.Client{}
	httpResp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != 204 {
		return fmt.Errorf("SendWatchMinute returned HTTP code: %d", httpResp.StatusCode)
	}
	
	return nil
}
