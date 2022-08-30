package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Adeithe/go-twitch"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"time"
	"ttv-cli/internal/pkg/twitch/gql/operation/redeemcustomreward"
	"ttv-cli/internal/pkg/twitch/pubsub/communitypointschannel"
)

type initialRewards []list.Item
type updatedReward communitypointschannel.UpdatedReward
type tick int

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initialRewards:
		cmd := m.list.SetItems(msg)
		return m, tea.Batch(cmd, m.processUpdates, m.tick())

	case updatedReward:
		// Reward has been paused or disabled, remove it from the list
		if msg.IsPaused || !msg.IsEnabled {
			for index, listItem := range m.list.Items() {
				if listItem.(*item).RewardId == msg.Id {
					m.list.RemoveItem(index)
				}
			}
			delete(m.itemsById, msg.Id)
		} else if item, ok := m.itemsById[msg.Id]; ok {
			item.CooldownExpiresAt = msg.CooldownExpiresAt
		}
		return m, m.processUpdates

	case tick:
		cmd := m.list.SetItems(m.list.Items()) // Force re-render
		return m, tea.Batch(cmd, m.tick())

	case tea.KeyMsg:
		if msg.String() == "enter" {
			selected := m.list.SelectedItem().(*item)
			// Don't attempt to redeem if the reward is still on cooldown
			if selected.GetExpiry().Seconds() > 0 {
				cmd := m.list.NewStatusMessage("Out of stock!")
				return m, cmd
			}
			m.redeemReward(selected)
			cmd := m.list.NewStatusMessage(fmt.Sprintf("Redeemed '%s'", selected.Title_))
			return m, cmd
		}

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

func (m Model) subscribeToRewards(ctx context.Context) {
	p := twitch.PubSub()
	err := p.Listen("community-points-channel-v1", m.twitchChannel.Id)
	if err != nil {
		log.Fatalln(err)
	}

	defer p.Close()

	handleUpdate := func(_ int, _ string, data []byte) {
		response := communitypointschannel.Response{}
		if err := json.Unmarshal(data, &response); err != nil {
			log.Fatalln(err)
		}
		m.rewardsUpdateChannel <- response
	}

	p.OnShardMessage(handleUpdate)

	<-ctx.Done()
}

func (m Model) processUpdates() tea.Msg {
	for update := range m.rewardsUpdateChannel {
		if update.Type == "custom-reward-updated" {
			return updatedReward(update.Data.UpdatedReward)
		}
	}
	return nil // Unreachable
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return tick(0)
	})
}

func (m Model) redeemReward(i *item) {

	input := redeemcustomreward.Input{
		ChannelID: m.twitchChannel.Id,
		Cost:      i.Cost,
		Prompt:    i.Prompt,
		RewardID:  i.RewardId,
		Title:     i.Title_,
	}

	if i.IsUserInputRequired {
		input.TextInput = ":)" // FIXME
	}

	_, err := redeemcustomreward.Redeem(input, m.authToken)
	if err != nil {
		log.Fatalln(err)
	}
}
