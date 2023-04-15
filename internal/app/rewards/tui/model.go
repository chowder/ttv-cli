package tui

import (
	"fmt"
	"github.com/Adeithe/go-twitch/pubsub"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"sort"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelpointscontext"
	"ttv-cli/internal/pkg/twitch/gql/query/channel"
	"ttv-cli/internal/pkg/twitch/pubsub/communitypointschannel"
	"ttv-cli/internal/pkg/twitch/pubsub/communitypointsuser"
)

type Model struct {
	twitchChannel        channel.Channel
	config               config.Config
	list                 list.Model
	itemsById            map[string]*item
	rewardsUpdateChannel chan communitypointschannel.Response
	pointsUpdateChannel  chan communitypointsuser.PointsSpentData
	notificationsChannel chan string
	pubsubClient         *pubsub.Client
}

func NewModel(pubsubClient *pubsub.Client, config config.Config, streamer string) Model {
	c, err := channel.GetChannel(config, streamer)
	if err != nil {
		log.Fatalf("Failed to get channel information for '%s' - %s", streamer, err)
	}
	if len(c.Id) == 0 {
		log.Fatalf("Could not find channel for '%s'\n", streamer)
	}

	m := Model{
		twitchChannel:        c,
		config:               config,
		list:                 list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		itemsById:            make(map[string]*item),
		rewardsUpdateChannel: make(chan communitypointschannel.Response, 8),
		pointsUpdateChannel:  make(chan communitypointsuser.PointsSpentData, 8),
		notificationsChannel: make(chan string, 8),
	}

	channelPointsContext, err := channelpointscontext.Get(m.config, c.Name)
	if err != nil {
		log.Fatalf("Could not fetch channel points context: %s", err)
	}

	balance := channelPointsContext.Data.Community.Channel.Self.CommunityPoints.Balance

	m.list.Title = fmt.Sprintf("%s's Rewards (%d points)", m.twitchChannel.DisplayName, balance)

	m.pubsubClient = pubsubClient

	return m
}

func (m Model) Init() tea.Cmd {
	go m.subscribeToRewards()
	go m.subscribeToPoints()
	return tea.Batch(m.getInitialRewards, m.processPointsUpdates, m.processNotificationUpdates, m.tick())
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
