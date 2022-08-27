package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"strings"
	"ttv-cli/internals/app/ttv-rewards/tui"
	"ttv-cli/internals/pkg/config"
	"ttv-cli/internals/pkg/twitch/login"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Provide a streamer whose channel points rewards you want to view")
		os.Exit(1)
	}

	c := config.CreateOrRead()
	if len(c.AuthToken) == 0 {
		c.AuthToken = login.GetAccessToken("", "")
		c.Save()
	}

	s := strings.ToLower(os.Args[1])
	m := tui.NewModel(s, c.AuthToken)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatalln("Error running program: ", err)
	}
}
