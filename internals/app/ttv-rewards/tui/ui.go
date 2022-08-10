package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"log"
	"sort"
	"time"
	"ttv-tools/internals/pkg/twitch/gql"
	"ttv-tools/internals/pkg/twitch/pubsub"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var orange = color.New(color.FgHiRed).SprintFunc()

type item struct {
	Title_            string
	Desc              string
	CooldownExpiresAt string
}

func (i item) Title() string {
	if i.CooldownExpiresAt == "" {
		return i.Title_
	}

	expireTime, err := time.Parse(time.RFC3339, i.CooldownExpiresAt)
	if err != nil {
		log.Fatalln(err)
	}
	expiresIn := time.Until(expireTime).Truncate(time.Second)
	if expiresIn.Seconds() <= 0 {
		return i.Title_
	}

	return orange(fmt.Sprintf("%s - %s", i.Title_, expiresIn.String()))
}

func (i item) Description() string { return i.Desc }
func (i item) FilterValue() string { return i.Title_ }

type Model struct {
	twitchChannel        gql.Channel
	list                 list.Model
	itemsById            map[string]*item
	rewardsUpdateChannel chan pubsub.CommunityPointsChannelResponse
}

func NewModel(streamer string) Model {
	m := Model{
		twitchChannel:        gql.GetChannel(streamer),
		list:                 list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		itemsById:            make(map[string]*item),
		rewardsUpdateChannel: make(chan pubsub.CommunityPointsChannelResponse),
	}
	m.list.Title = "Rewards"
	return m
}

func (m Model) Init() tea.Cmd {
	ctx := context.Background() // TODO: Close this context on app exit
	go func() {
		if err := pubsub.CommunityPointsChannel(ctx, m.twitchChannel.Id, m.rewardsUpdateChannel); err != nil {
			log.Fatalln(err)
		}
	}()
	return tea.Batch(m.getRewards, m.tick())
}

type initialRewards []list.Item
type updatedReward pubsub.UpdatedReward
type tick int

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initialRewards:
		var cmd tea.Cmd
		cmd = m.list.SetItems(msg)
		return m, tea.Batch(cmd, m.processUpdates, m.tick())
	case updatedReward:
		item := m.itemsById[msg.Id]
		item.CooldownExpiresAt = msg.CooldownExpiresAt
		m.list.Title = fmt.Sprintf("Latest redemption: %s", msg.Title)
		return m, m.processUpdates
	case tick:
		var cmd tea.Cmd
		cmd = m.list.SetItems(m.list.Items())
		return m, tea.Batch(cmd, m.tick(), m.processUpdates)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func (m Model) getRewards() tea.Msg {
	rewards := m.twitchChannel.CommunityPointsSettings.CustomRewards
	sort.Slice(rewards[:], func(l, r int) bool {
		return rewards[l].Cost < rewards[r].Cost
	})

	items := make(initialRewards, 0)
	for _, reward := range rewards {
		title := fmt.Sprintf("%s (%d points)", reward.Title, reward.Cost)
		description := reward.Prompt
		item := item{Title_: title, Desc: description, CooldownExpiresAt: reward.CooldownExpiresAt}
		m.itemsById[reward.Id] = &item
		items = append(items, item)
	}

	return items
}

func (m Model) processUpdates() tea.Msg {
	update := <-m.rewardsUpdateChannel
	if update.Type == "custom-reward-updated" {
		return updatedReward(update.Data.UpdatedReward)
	}
	return nil
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return tick(0)
	})
}
