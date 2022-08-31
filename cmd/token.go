package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"ttv-cli/internal/pkg/twitch/login"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate a Twitch OAuth token for your account",
	Run: func(cmd *cobra.Command, args []string) {
		authToken, err := login.GetAccessToken(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_PASSWORD"))
		if err != nil {
			log.Fatalf("Could not fetch Twitch access token: %s\n", err)
		}
		fmt.Println(authToken)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
