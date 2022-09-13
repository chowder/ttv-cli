package miner

import (
	"context"
	"encoding/json"
	"github.com/Adeithe/go-twitch/pubsub"
	"log"
	"strings"
	"ttv-cli/internal/pkg/twitch/pubsub/communitypointsuser"
)

type pointsUpdate communitypointsuser.PointsEarnedData

func (m Miner) subscribePoints(ctx context.Context) error {
	log.Println("Subscribing to points updates...")

	c, err := getPointsChannel(m.pubsubClient, m.AuthToken, m.UserId)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case u := <-c:
				m.eventBus.Publish(pointsUpdateTopic, u)
			}
		}
	}()

	_ = m.eventBus.Subscribe(pointsUpdateTopic, listenToPointsUpdates)
	go func() {
		<-ctx.Done()
		_ = m.eventBus.Unsubscribe(pointsUpdateTopic, listenToPointsUpdates)
	}()

	return nil
}

func listenToPointsUpdates(update pointsUpdate) {
	log.Printf("+%d points in channel %s, reason: %s, balance: %d\n", update.PointGain.TotalPoints, update.ChannelId, update.PointGain.ReasonCode, update.Balance.Balance)
}

func getPointsChannel(client *pubsub.Client, authToken string, userId string) (<-chan pointsUpdate, error) {
	const topic = "community-points-user-v1"

	c := make(chan pointsUpdate)

	handleUpdate := func(s int, t string, data []byte) {
		if strings.HasPrefix(t, topic) {
			var resp communitypointsuser.Response
			if err := json.Unmarshal(data, &resp); err != nil {
				log.Printf("could not unmarshal response: %s, error %s\n", string(data), err)
				return
			}

			if resp.Type == "points-earned" {
				var p communitypointsuser.PointsEarnedData
				if err := json.Unmarshal(resp.Data, &p); err != nil {
					log.Printf("could not unmarshal response: %s, error %s\n", string(resp.Data), err)
					return
				}

				c <- pointsUpdate(p)
			}
		}
	}

	client.OnShardMessage(handleUpdate)
	err := client.ListenWithAuth(authToken, topic, userId)
	if err != nil {
		return nil, err
	}

	return c, nil
}
