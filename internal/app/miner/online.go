package miner

import (
	"context"
	"fmt"
	"github.com/asaskevich/EventBus"
	"golang.org/x/exp/slices"
	"log"
	"time"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
	"ttv-cli/internal/pkg/utils"
)

type StreamDetails struct {
	login  string
	stream *users.Stream
}

func (m Miner) subscribeStreamStatus(ctx context.Context) error {
	log.Println("Subscribing to stream start events...")

	c, err := getStreamStartChannel(m.client, ctx)
	if err != nil {
		return fmt.Errorf("could not create stream start channel: %w", err)
	}

	activeStreamers := make([]string, 0)

	go func() {
		for details := range c {
			if details.stream == nil && slices.Contains(activeStreamers, details.login) {
				m.eventBus.Publish(streamEndTopic, details)
				activeStreamers = utils.Remove(activeStreamers, details.login)
			} else if details.stream != nil && !slices.Contains(activeStreamers, details.login) {
				m.eventBus.Publish(streamStartTopic, details)
				activeStreamers = append(activeStreamers, details.login)
			}
		}
	}()

	if err := registerStreamStartHandlers(ctx, m.eventBus); err != nil {
		return err
	}

	return nil
}

func registerStreamStartHandlers(ctx context.Context, eventBus EventBus.Bus) error {
	streamStartHandler := func(stream StreamDetails) {
		log.Println("Streamer came online: ", stream.login)
	}

	streamEndHandler := func(stream StreamDetails) {
		log.Println("Streamer went offline: ", stream.login)
	}

	if err := eventBus.Subscribe(streamStartTopic, streamStartHandler); err != nil {
		return err
	}
	if err := eventBus.Subscribe(streamEndTopic, streamEndHandler); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		_ = eventBus.Unsubscribe(streamStartTopic, streamStartHandler)
		_ = eventBus.Unsubscribe(streamEndTopic, streamEndHandler)
	}()

	return nil
}

func getStreamStartChannel(config config.Config, ctx context.Context) (<-chan StreamDetails, error) {
	follows, err := channelfollows.Get(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get followed channels: %w", err)
	}

	logins := make([]string, len(follows))
	for i, f := range follows {
		logins[i] = f.Login
	}

	c := make(chan StreamDetails)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c)
				return
			case <-time.After(30 * time.Second):
				us, err := users.GetUsers(config, logins)
				if err != nil {
					log.Println("could not get updates: ", err)
				}
				for _, user := range us {
					c <- StreamDetails{
						login:  user.Login,
						stream: user.Stream,
					}
				}
			}
		}
	}()

	return c, nil
}
