package utils

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"time"
)

import (
	"ttv-live/srcs/twitch"
)

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute

	if h == 0 {
		return fmt.Sprintf("%02dm", m)
	}

	return fmt.Sprintf("%1dh%02dm", h, m)
}

func DisplayUserLive(user twitch.User) {
	green := color.New(color.FgHiGreen).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()

	streamer := fmt.Sprintf("%-10s", user.DisplayName)
	directory := user.Stream.Game.DisplayName

	startTime, err := time.Parse(time.RFC3339, user.Stream.CreatedAt)
	if err != nil {
		log.Fatalln(err)
	}

	duration := time.Since(startTime)
	duration = duration.Truncate(time.Second)

	fmt.Println("  \u2022", green(streamer), ":", green(directory), yellow("(", FmtDuration(duration), ")"))
}

func DisplayUserOffline(user twitch.User) {
	red := color.New(color.FgHiRed).SprintFunc()
	streamer := fmt.Sprintf("%-10s", user.DisplayName)

	fmt.Println("  \u2022", red(streamer), ":", red("(offline)"))
}
