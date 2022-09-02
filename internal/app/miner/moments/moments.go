package moments

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Adeithe/go-twitch/pubsub"
	"log"
	"strings"
	"time"
	"ttv-cli/internal/pkg/twitch/gql/operation/communitymomentcalloutclaim"
	"ttv-cli/internal/pkg/twitch/pubsub/communitymomentschannel"
)

const topic = "community-moments-channel-v1"

func MineMoments(c *pubsub.Client, streamerByIds map[string]string, authToken string) error {
	for id, s := range streamerByIds {
		log.Printf("Listening to topic: '%s' for streamer: '%s' (%s)\n", topic, s, id)
		if err := c.ListenWithAuth(authToken, topic, id); err != nil {
			msg := fmt.Sprintf("Failed to listen to topic: '%s' for streamer: '%s' (%s) - %v", topic, s, id, err)
			return errors.New(msg)
		}
		time.Sleep(time.Second)
	}

	handleUpdate := func(shard int, topic string, data []byte) {
		fmt.Printf("Shard #%d > %s %s\n", shard, topic, strings.TrimSpace(string(data)))

		var resp communitymomentschannel.Response
		if err := json.Unmarshal(data, &resp); err != nil {
			log.Printf("could not unmarshal response: %s, error %s\n", string(data), err)
			return
		}

		var message communitymomentschannel.Message
		if err := json.Unmarshal([]byte(resp.Data.Message), &message); err != nil {
			log.Printf("could not unmarshal response message: %s, error: %s\n", resp.Data.Message, err)
			return
		}

		if message.Type == "active" {
			log.Printf("Active moment received, data: %s\n", message.Data)
			var data communitymomentschannel.ActiveMomentData
			if err := json.Unmarshal(message.Data, &data); err != nil {
				log.Printf("could not unmarshal response message data: %s, error: %s\n", message.Data, err)
				return
			}

			log.Printf("Attempting to redeem moment ID: '%s'\n", data.MomentId)
			err := communitymomentcalloutclaim.Claim(data.MomentId, authToken)
			if err != nil {
				log.Printf("could not claim moment: %s, error: %s\n", data.MomentId, err)
			}
		}
	}

	c.OnShardMessage(handleUpdate)

	return nil
}
