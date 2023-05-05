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

	m := miner.New(c)
	m.Start()
}

func init() {
	rootCmd.AddCommand(minerCmd)
}
