package miner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Adeithe/go-twitch/pubsub"
	"github.com/asaskevich/EventBus"
	"log"
	"strings"
	"time"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
	"ttv-cli/internal/pkg/twitch/gql/operation/communitymomentcalloutclaim"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
	"ttv-cli/internal/pkg/twitch/pubsub/communitymomentschannel"
	"ttv-cli/internal/pkg/utils"
)

type Moment communitymomentschannel.Response

func (m Miner) subscribeMoments(ctx context.Context) error {
	log.Println("Subscribing to moments...")

	c, err := getMomentsChannel(m.pubsubClient, m.AuthToken)
	if err != nil {
		return err
	}

	go func() {
		for moment := range c {
			m.eventBus.Publish(momentsTopic, moment)
		}
	}()

	if err := registerMomentsHandlers(ctx, m.AuthToken, m.eventBus); err != nil {
		return err
	}

	return nil
}

func registerMomentsHandlers(ctx context.Context, authToken string, eventBus EventBus.Bus) error {
	handler := func(moment Moment) {
		momentId := moment.Data.MomentId
		if len(momentId) > 0 {
			log.Printf("Attempting to redeem moment ID: '%s'\n", momentId)
			err := communitymomentcalloutclaim.Claim(momentId, authToken)
			if err != nil {
				log.Printf("could not claim moment: %s, error: %s\n", momentId, err)
			}
		}
	}

	go func() {
		<-ctx.Done()
		_ = eventBus.Unsubscribe(momentsTopic, handler)
	}()

	return eventBus.Subscribe(momentsTopic, handler)
}

func getMomentsChannel(client *pubsub.Client, authToken string) (<-chan Moment, error) {
	followsById, err := getFollowsById(authToken)
	if err != nil {
		return nil, err
	}

	const topic = "community-moments-channel-v1"

	success := make([]string, 0)
	for id, name := range followsById {
		if err := client.ListenWithAuth(authToken, topic, id); err != nil {
			msg := fmt.Sprintf("Failed to listen to topic: '%s' for streamer: '%s' (%s) - %v", topic, name, id, err)
			log.Println(msg)
		}
		success = append(success, name)
		time.Sleep(utils.GetRandomDuration(1, 2))
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

	client.OnShardMessage(handleUpdate)

	log.Printf("Mining Moments for users: %v\n", success)

	return c, nil
}

func getFollowsById(authToken string) (map[string]string, error) {
	follows, err := channelfollows.Get(authToken)
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
		followsByIds[userInfo.Id] = userInfo.Login
	}

	return followsByIds, nil
}
