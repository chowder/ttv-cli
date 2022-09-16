package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"ttv-cli/internal/app/miner"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch"
)

var minerCmd = &cobra.Command{
	Use:   "miner",
	Short: "Mines channel points, moments, and drops",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	c, err := config.CreateOrRead()
	if err != nil {
		log.Fatalln("could not read config file: ", err)
	}

	client := twitch.NewClient(c.AuthToken)

	m := miner.New(client, c.TokenDetails.UserId)
	m.Start()
}

func init() {
	rootCmd.AddCommand(minerCmd)
}
