package miner

import (
	"context"
	"fmt"
	"log"
	"time"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
)

const streamStartTopic = "stream_start"

func (m Miner) listenStreamStart(ctx context.Context) error {
	log.Println("listening to stream start events...")

	handler := func(name string) {
		log.Println("Streamer came online: ", name)
	}

	go func() {
		<-ctx.Done()
		_ = m.eventBus.Unsubscribe(streamStartTopic, handler)
	}()

	return m.eventBus.Subscribe(streamStartTopic, handler)
}

func (m Miner) subscribeStreamStart(ctx context.Context) error {
	log.Println("subscribing to stream start events...")

	c, err := getStreamStartChannel(ctx, m.AuthToken)
	if err != nil {
		return fmt.Errorf("could not create stream start channel: %w", err)
	}

	go func() {
		for {
			select {
			case stream := <-c:
				m.eventBus.Publish(streamStartTopic, stream.DisplayName)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func getStreamStartChannel(ctx context.Context, authToken string) (<-chan users.User, error) {
	follows, err := channelfollows.Get(authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get followed channels: %w", err)
	}

	streamStartTimesByLogins := make(map[string]time.Time)
	logins := make([]string, len(follows))
	for i, f := range follows {
		streamStartTimesByLogins[f.Login] = time.Time{}
		logins[i] = f.Login
	}

	c := make(chan users.User)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c)
				return
			case <-time.After(30 * time.Second):
				us, err := users.GetUsers(logins)
				if err != nil {
					log.Println("could not get updates: ", err)
				}
				for _, user := range us {
					previousStreamStartTime := streamStartTimesByLogins[user.Login]
					if user.Stream != nil && !user.Stream.CreatedAt.Equal(previousStreamStartTime) {
						c <- user
						streamStartTimesByLogins[user.Login] = user.Stream.CreatedAt
					}
				}
			}
		}
	}()

	return c, nil
}
