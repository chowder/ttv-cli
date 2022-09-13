package live

import (
	"fmt"
	"github.com/fatih/color"
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

func DisplayUserLive(user users.User, width int) {
	green := color.New(color.FgHiGreen).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()

	streamer := fmt.Sprintf("%-*s", width, user.DisplayName)
	directory := user.Stream.Game.DisplayName

	duration := time.Since(user.Stream.CreatedAt)
	duration = duration.Truncate(time.Second)

	fmt.Println("  \u2022", green(streamer), ":", green(directory), yellow("(", FmtDuration(duration), ")"))
}

func DisplayUserOffline(user users.User, width int) {
	red := color.New(color.FgHiRed).SprintFunc()
	streamer := fmt.Sprintf("%-*s", width, user.DisplayName)

	fmt.Println("  \u2022", red(streamer), ":", red("(offline)"))
}
