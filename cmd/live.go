package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math"
	"ttv-cli/internal/app/live"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
)

var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "View which of your follows are currently live",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Load()
		if err != nil {
			c, err = config.Create()
			if err != nil {
				log.Fatalln(err)
			}
			err = c.Save()
			if err != nil {
				log.Fatalln(err)
			}
		}

		f, err := channelfollows.Get(c)
		if err != nil {
			log.Fatalf("Error fetching followed channels: %s\n", err)
		}

		s := make([]string, 0)
		for _, f := range f {
			s = append(s, f.Login)
		}

		// Get all streamers from Twitch API
		streamers, err := users.GetUsers(c, s)
		if err != nil {
			log.Fatalf("Could not get channel information - %s\n", err)
		}

		// Filter between live and offline streamers
		online := make([]users.User, 0)
		offline := make([]users.User, 0)

		for i, user := range streamers {
			if len(user.Id) == 0 {
				fmt.Printf("Could not find channel information for '%s'\n", s[i])
				continue
			} else if user.Stream != nil && !user.Stream.CreatedAt.IsZero() {
				online = append(online, user)
			} else {
				offline = append(offline, user)
			}
		}

		width := 0
		for _, user := range online {
			width = int(math.Max(float64(width), float64(len(user.DisplayName))))
		}
		if showOffline {
			for _, user := range offline {
				width = int(math.Max(float64(width), float64(len(user.DisplayName))))
			}
		}

		// Display online streamers first, then offline ones
		for _, user := range online {
			live.DisplayUserLive(user, width)
		}

		if showOffline {
			for _, user := range offline {
				live.DisplayUserOffline(user, width)
			}
		}
	},
}

var showOffline bool

func init() {
	rootCmd.AddCommand(liveCmd)
	liveCmd.Flags().BoolVarP(&showOffline, "show-offline", "a", false, "Toggle to display streamers who are offline")
}
