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
	redeem "ttv-cli/internals/pkg/twitch/gql/redeem_custom_reward"
	"ttv-cli/internals/pkg/twitch/pubsub"
)

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
		if msg.IsPaused || !msg.IsEnabled {
			delete(m.itemsById, msg.Id)
		} else {
			item.CooldownExpiresAt = msg.CooldownExpiresAt
		}
		return m, m.processUpdates
	case tick:
		cmd := m.list.SetItems(m.list.Items()) // Force re-render
		return m, tea.Batch(cmd, m.tick(), m.processUpdates)
	case tea.KeyMsg:
		if msg.String() == "enter" {
			selected := m.list.SelectedItem().(*item)
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
		response := pubsub.CommunityPointsChannelResponse{}
		if err := json.Unmarshal(data, &response); err != nil {
			log.Fatalln(err)
		}
		m.rewardsUpdateChannel <- response
	}

	p.OnShardMessage(handleUpdate)

	<-ctx.Done()
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

func (m Model) redeemReward(i *item) {
	_, err := redeem.RedeemCustomReward(m.twitchChannel.Id, i.Cost, i.Prompt, i.RewardId, i.Title_, m.authToken)
	if err != nil {
		log.Fatalln(err)
	}
}
