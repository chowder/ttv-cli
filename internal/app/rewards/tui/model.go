package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"sort"
	"ttv-cli/internal/pkg/twitch/gql/query/channel"
	"ttv-cli/internal/pkg/twitch/pubsub"
)

type Model struct {
	twitchChannel        channel.Channel
	authToken            string
	list                 list.Model
	itemsById            map[string]*item
	rewardsUpdateChannel chan pubsub.CommunityPointsChannelResponse
}

func NewModel(streamer string, authToken string) Model {
	//listModel := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	m := Model{
		twitchChannel:        channel.GetChannel(streamer),
		authToken:            authToken,
		list:                 list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		itemsById:            make(map[string]*item),
		rewardsUpdateChannel: make(chan pubsub.CommunityPointsChannelResponse),
	}
	m.list.Title = fmt.Sprintf("%s's Rewards", m.twitchChannel.DisplayName)
	return m
}

func (m Model) Init() tea.Cmd {
	ctx := context.Background() // TODO: Close this context on app exit
	go m.subscribeToRewards(ctx)
	return tea.Batch(m.getInitialRewards, m.tick())
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func (m Model) getInitialRewards() tea.Msg {
	rewards := m.twitchChannel.CommunityPointsSettings.CustomRewards
	sort.Slice(rewards[:], func(l, r int) bool {
		return rewards[l].Cost < rewards[r].Cost
	})

	items := make(initialRewards, 0)
	for _, reward := range rewards {
		if reward.IsPaused || !reward.IsEnabled {
			continue
		}
		item := item{
			Title_:              reward.Title,
			Prompt:              reward.Prompt,
			Cost:                reward.Cost,
			RewardId:            reward.Id,
			CooldownExpiresAt:   reward.CooldownExpiresAt,
			IsUserInputRequired: reward.IsUserInputRequired,
		}
		m.itemsById[reward.Id] = &item
		items = append(items, &item)
	}

	return items
}
