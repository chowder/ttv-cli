package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"strings"
	"ttv-tools/internals/app/ttv-rewards/tui"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Provide a streamer whose channel points rewards you want to view")
		os.Exit(1)
	}

	s := strings.ToLower(os.Args[1])
	m := tui.NewModel(s)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatalln("Error running program: ", err)
	}
}
