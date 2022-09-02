package tui

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"time"
	"ttv-cli/internal/pkg/twitch/gql/operation/redeemcustomreward"
	"ttv-cli/internal/pkg/twitch/pubsub/communitypointschannel"
	"ttv-cli/internal/pkg/twitch/pubsub/communitypointsuser"
)

type initialRewards []list.Item
type updatedReward communitypointschannel.UpdatedReward
type tick int
type newBalance int

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initialRewards:
		cmd := m.list.SetItems(msg)
		return m, tea.Batch(cmd, m.processRewardUpdates, m.tick())

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
		return m, m.processRewardUpdates

	case newBalance:
		m.list.Title = fmt.Sprintf("%s's Rewards (%d points)", m.twitchChannel.DisplayName, msg)
		return m, m.processPointsUpdates

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
			err := m.redeemReward(selected)
			if err != nil {
				cmd := m.list.NewStatusMessage("Could not redeem!")
				return m, cmd
			}
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

func (m Model) subscribeToRewards() {
	err := m.pubsubClient.Listen("community-points-channel-v1", m.twitchChannel.Id)
	if err != nil {
		log.Fatalf("Could not subscribe to community-points-channel-v1: %s\n", err)
	}

	subscribedTopic := "community-points-channel-v1." + m.twitchChannel.Id
	handleUpdate := func(_ int, topic string, data []byte) {
		if topic == subscribedTopic {
			var response communitypointschannel.Response
			if err := json.Unmarshal(data, &response); err != nil {
				log.Fatalln(err)
			}
			m.rewardsUpdateChannel <- response
		}
	}

	m.pubsubClient.OnShardMessage(handleUpdate)
}

func (m Model) subscribeToPoints() {
	userId := m.config.TokenDetails.UserId

	err := m.pubsubClient.ListenWithAuth(m.config.AuthToken, "community-points-user-v1", userId)
	if err != nil {
		log.Fatalf("Could not subscribe to community-points-user-v1.%s: %s\n", m.twitchChannel.Id, err)
	}

	subscribedTopic := "community-points-user-v1." + userId
	handleUpdate := func(_ int, topic string, data []byte) {
		if topic == subscribedTopic {
			var response communitypointsuser.Response
			if err := json.Unmarshal(data, &response); err != nil {
				log.Fatalf("Failed to unmarshal response: %s, error: %s\n", data, err)
			}

			if response.Type == "points-earned" || response.Type == "points-spent" {
				m.pointsUpdateChannel <- response
			}
		}
	}

	m.pubsubClient.OnShardMessage(handleUpdate)
}

func (m Model) processRewardUpdates() tea.Msg {
	for update := range m.rewardsUpdateChannel {
		if update.Type == "custom-reward-updated" {
			return updatedReward(update.Data.UpdatedReward)
		}
	}
	return nil // Unreachable
}

func (m Model) processPointsUpdates() tea.Msg {
	for update := range m.pointsUpdateChannel {
		if update.Type == "points-spent" {
			var data communitypointsuser.PointsSpentData
			if err := json.Unmarshal(update.Data, &data); err != nil {
				log.Fatalf("Could not process point update: %s, error: %s\n", update, err)
			}
			return newBalance(data.Balance.Balance)
		} else if update.Type == "points-earned" {
			var data communitypointsuser.PointsEarnedData
			if err := json.Unmarshal(update.Data, &data); err != nil {
				log.Fatalf("Could not process point update: %s, error: %s\n", update, err)
			}
			return newBalance(data.Balance.Balance)
		}
	}
	return nil
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return tick(0)
	})
}

func (m Model) redeemReward(i *item) error {

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

	_, err := redeemcustomreward.Redeem(input, m.config.AuthToken)
	if err != nil {
		return fmt.Errorf("could not redeem reward: %w", err)
	}

	return nil
}
