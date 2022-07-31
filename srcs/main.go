package main

import (
	"log"
)

import (
	"ttv-live/srcs/config"
	"ttv-live/srcs/twitch"
	"ttv-live/srcs/utils"
)

func main() {
	c := config.CreateOrReadFromFile(config.GetDefaultConfigFile())
	if len(c.Streamers) == 0 {
		log.Fatalf("No streamers specified in config, please populate '%s' with a list of streamers\n", config.GetDefaultConfigFile())
	}

	// Get all streamers from Twitch API
	streamers := twitch.GetAllStreamers(c.Streamers)

	// Filter between live and offline streamers
	online := make([]twitch.User, 0)
	offline := make([]twitch.User, 0)

	for _, user := range streamers {
		if user.Stream.CreatedAt != "" {
			online = append(online, user)
		} else {
			offline = append(offline, user)
		}
	}

	// Display to terminal
	for _, user := range online {
		utils.DisplayUserLive(user)
	}

	for _, user := range offline {
		utils.DisplayUserOffline(user)
	}
}
