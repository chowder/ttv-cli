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
			log.Println(err)
		}

		if len(resp.Data.MomentId) > 0 {
			log.Printf("Attempting to redeem moment ID: '%s'\n", resp.Data.MomentId)
			err := communitymomentcalloutclaim.Claim(resp.Data.MomentId, authToken)
			if err != nil {
				log.Println(err)
			}
		}
	}

	c.OnShardMessage(handleUpdate)

	return nil
}
