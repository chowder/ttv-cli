package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"ttv-cli/internal/pkg/twitch/auth"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate a Twitch OAuth token for your account",
	Run: func(cmd *cobra.Command, args []string) {
		authToken, err := auth.GetAccessToken()
		if err != nil {
			log.Fatalf("Could not fetch Twitch access token: %s\n", err)
		}
		fmt.Println(authToken)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
