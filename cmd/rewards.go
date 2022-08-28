package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"ttv-cli/internal/app/rewards/tui"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/login"
)

var rewardsCmd = &cobra.Command{
	Use:        "rewards STREAMER_NAME",
	Short:      "View and redeem Twitch channel point rewards",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"streamer"},
	Run: func(cmd *cobra.Command, args []string) {
		c := config.CreateOrRead()
		if len(c.AuthToken) == 0 {
			authToken, err := login.GetAccessToken("", "")
			if err != nil {
				log.Fatalln(err)
			}

			c.AuthToken = authToken
			c.Save()
		}

		s := strings.ToLower(args[0])
		m := tui.NewModel(s, c.AuthToken)

		p := tea.NewProgram(m, tea.WithAltScreen())

		if err := p.Start(); err != nil {
			log.Fatalln("Error running program: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rewardsCmd)
}
