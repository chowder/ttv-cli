package cmd

import (
	"github.com/spf13/cobra"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
	"ttv-cli/internal/pkg/utils"
)

var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "View which of your follows are currently live",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.CreateOrRead()

		f := channelfollows.GetChannelFollows(c.AuthToken)
		s := make([]string, 0)
		for _, f := range f {
			s = append(s, f.Login)
		}

		// Get all streamers from Twitch API
		streamers := users.GetUsers(s)

		// Filter between live and offline streamers
		online := make([]users.User, 0)
		offline := make([]users.User, 0)

		for _, user := range streamers {
			if user.Stream.CreatedAt != "" {
				online = append(online, user)
			} else {
				offline = append(offline, user)
			}
		}

		// Display online streamers first, then offline ones
		for _, user := range online {
			utils.DisplayUserLive(user)
		}

		if showOffline {
			for _, user := range offline {
				utils.DisplayUserOffline(user)
			}
		}
	},
}

var showOffline bool

func init() {
	rootCmd.AddCommand(liveCmd)
	showOffline = *liveCmd.Flags().BoolP("show-offline", "a", false, "Toggle to display streamers who are offline")
}
