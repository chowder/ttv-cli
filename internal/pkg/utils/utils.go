package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"log"
	"time"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
)

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute

	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}

	return fmt.Sprintf("%1dh%02dm", h, m)
}

func DisplayUserLive(user users.User) {
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

func DisplayUserOffline(user users.User) {
	red := color.New(color.FgHiRed).SprintFunc()
	streamer := fmt.Sprintf("%-10s", user.DisplayName)

	fmt.Println("  \u2022", red(streamer), ":", red("(offline)"))
}

func TokenHex(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln("Error reading random token: ", err)
	}
	return hex.EncodeToString(b)
}
