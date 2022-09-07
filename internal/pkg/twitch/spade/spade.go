package spade

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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
