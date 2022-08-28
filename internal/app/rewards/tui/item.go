package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"log"
	"time"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var orange = color.New(color.FgHiRed).SprintFunc()

type item struct {
	Title_              string
	Cost                int
	Prompt              string
	RewardId            string
	CooldownExpiresAt   string
	IsUserInputRequired bool
}

func (i item) Title() string {
	title := fmt.Sprintf("%s (%d points)", i.Title_, i.Cost)

	if i.CooldownExpiresAt == "" {
		return title
	}

	expiresIn := i.GetExpiry()
	if expiresIn.Seconds() <= 0 {
		return title
	}

	return orange(fmt.Sprintf("%s - %s", title, expiresIn.String()))
}

func (i item) Description() string { return i.Prompt }
func (i item) FilterValue() string { return i.Title_ }

func (i item) GetExpiry() time.Duration {
	if i.CooldownExpiresAt == "" {
		return 0
	}

	expireTime, err := time.Parse(time.RFC3339, i.CooldownExpiresAt)
	if err != nil {
		log.Fatalln(err)
	}

	return time.Until(expireTime).Truncate(time.Second)
}
