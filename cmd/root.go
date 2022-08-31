package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var buildVersion = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ttv",
	Short:   "A collection of command line sub-programs for Twitch",
	Version: buildVersion,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
