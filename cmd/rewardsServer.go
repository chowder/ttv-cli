package cmd

import (
	"github.com/spf13/cobra"
	"ttv-cli/internal/app/rewards/server"
)

var (
	serverAddr string
)

var rewardsServerCmd = &cobra.Command{
	Use:   "rewards-server",
	Short: "Serve a webpage showing redemption statuses and cooldowns",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run(serverAddr)
	},
}

func init() {
	rootCmd.AddCommand(rewardsServerCmd)

	rewardsServerCmd.Flags().StringVarP(&serverAddr, "addr", "a", "localhost:3000", "address to serve the server on")
}
