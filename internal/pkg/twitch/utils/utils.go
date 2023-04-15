package utils

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"regexp"
)

var clientVersionPattern = regexp.MustCompile("window\\.__twilightBuildID=\"(.*?)\"")

func GetClientVersion() (string, error) {
	client := resty.New()
	resp, err := client.R().
		Get("https://www.twitch.tv/")

	if err != nil {
		return "", err
	}

	matches := clientVersionPattern.FindSubmatch(resp.Body())
	if len(matches) < 2 {
		println(resp.String())
		return "", errors.New("could not find client version")
	}

	return string(matches[1]), nil
}
