package config

import (
	"errors"
	"math/rand"
	"regexp"
	"ttv-cli/internal/pkg/twitch"
	"ttv-cli/internal/pkg/utils"
)

var clientVersionPattern = regexp.MustCompile("window\\.__twilightBuildID=\"(.*?)\"")

func getClientVersion() (string, error) {
	resp, err := utils.HttpGet(twitch.HomeUrl, nil)
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
