package miner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
	"ttv-cli/internal/pkg/twitch/gql/operation/communitymomentcalloutclaim"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
	"ttv-cli/internal/pkg/twitch/pubsub/communitymomentschannel"
)

const momentsTopic = "moment"

type Moment communitymomentschannel.Response

func (m Miner) listenMoments(ctx context.Context) error {
	handler := func(moment Moment) {
		momentId := moment.Data.MomentId
		if len(momentId) > 0 {
			log.Printf("Attempting to redeem moment ID: '%s'\n", momentId)
			err := communitymomentcalloutclaim.Claim(momentId, m.AuthToken)
			if err != nil {
				log.Printf("could not claim moment: %s, error: %s\n", momentId, err)
			}
		}
	}

	go func() {
		<-ctx.Done()
		_ = m.eventBus.Unsubscribe(momentsTopic, handler)
	}()

	return m.eventBus.Subscribe(momentsTopic, handler)
}

func (m Miner) subscribeMoments(ctx context.Context) error {
	c, err := m.getMomentsChannel(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case moment := <-c:
				m.eventBus.Publish(momentsTopic, moment)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (m Miner) getMomentsChannel(ctx context.Context) (<-chan Moment, error) {
	followsById, err := m.getFollowsById()
	if err != nil {
		return nil, err
	}

	const topic = "community-moments-channel-v1"

	for id, name := range followsById {
		log.Printf("Listening to topic: '%s' for streamer: '%s' (%s)\n", topic, name, id)
		if err := m.pubsubClient.ListenWithAuth(m.AuthToken, topic, id); err != nil {
			msg := fmt.Sprintf("Failed to listen to topic: '%s' for streamer: '%s' (%s) - %v", topic, name, id, err)
			log.Println(msg)
		}
		time.Sleep(time.Second) // TODO: Add jitter
	}

	c := make(chan Moment)

	handleUpdate := func(s int, t string, data []byte) {
		if strings.HasPrefix(t, topic) {
			var resp communitymomentschannel.Response
			if err := json.Unmarshal(data, &resp); err != nil {
				log.Printf("could not unmarshal response: %s, error %s\n", string(data), err)
				return
			}

			c <- Moment(resp)
		}
	}

	m.pubsubClient.OnShardMessage(handleUpdate)

	return c, nil
}

func (m Miner) getFollowsById() (map[string]string, error) {
	follows, err := channelfollows.Get(m.AuthToken)
	if err != nil {
		return nil, err
	}

	loginsOfFollows := make([]string, len(follows))
	for i, f := range follows {
		loginsOfFollows[i] = f.Login
	}

	userInfos, err := users.GetUsers(loginsOfFollows)
	if err != nil {
		return nil, err
	}

	followsByIds := make(map[string]string)
	for _, userInfo := range userInfos {
		followsByIds[userInfo.Id] = userInfo.DisplayName
	}

	return followsByIds, nil
}
