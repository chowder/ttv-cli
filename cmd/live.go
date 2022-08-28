package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql/users"
	"ttv-cli/internal/pkg/utils"
)

var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "View which streamers are currently live",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.CreateOrRead()
		if len(c.Streamers) == 0 {
			log.Fatalf("No streamers specified in config, please populate '%s' with a list of streamers\n", config.GetConfigFilePath())
		}

		// Get all streamers from Twitch API
		streamers := users.GetUsers(c.Streamers)

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

		// Display to terminal
		for _, user := range online {
			utils.DisplayUserLive(user)
		}

		for _, user := range offline {
			utils.DisplayUserOffline(user)
		}
	},
}

func init() {
	rootCmd.AddCommand(liveCmd)
}
