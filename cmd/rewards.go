package cmd

import (
	"github.com/Adeithe/go-twitch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"ttv-cli/internal/app/rewards/tui"
	"ttv-cli/internal/pkg/config"
)

var rewardsCmd = &cobra.Command{
	Use:        "rewards STREAMER_NAME",
	Short:      "View and redeem Twitch channel point rewards",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"streamer"},
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.CreateOrRead()
		if err != nil {
			log.Fatalf("Error reading config: %s\n", err)
		}

		s := strings.ToLower(args[0])

		pc := twitch.PubSub()
		defer pc.Close()

		m := tui.NewModel(pc, c, s)

		p := tea.NewProgram(m, tea.WithAltScreen())

		if err := p.Start(); err != nil {
			log.Fatalln("Error running program: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rewardsCmd)
}
