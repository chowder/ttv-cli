package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"ttv-cli/internal/app/miner"
	"ttv-cli/internal/pkg/config"
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
	m := miner.New(c.AuthToken)
	m.Start()
}

func init() {
	rootCmd.AddCommand(minerCmd)
}
