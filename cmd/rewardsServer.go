package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ttv-cli/internal/app/rewards/server"
)

var (
	serverHost string
	serverPort int
	streamer   string
)

var rewardsServerCmd = &cobra.Command{
	Use:   "rewards-server",
	Short: "Serve a webpage showing redemption statuses and cooldowns",
	Run: func(cmd *cobra.Command, args []string) {
		addr := fmt.Sprintf("%s:%d", serverHost, serverPort)
		server.Run(addr, streamer)
	},
}

func init() {
	rootCmd.AddCommand(rewardsServerCmd)

	rewardsServerCmd.Flags().StringVarP(&serverHost, "host", "", "localhost", "hostname to bind the server to")
	rewardsServerCmd.Flags().IntVarP(&serverPort, "port", "", 8080, "port to bind the server to")
	rewardsServerCmd.Flags().StringVarP(&streamer, "streamer", "", "", "streamer whose redemptions to track")
}
