package pubsub

import (
	"fmt"
	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/pubsub"
	"os"
	"time"
)

func RecordedPubSub(path string) (*pubsub.Client, *os.File, error) {
	c := twitch.PubSub()
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	handleUpdate := func(shard int, topic string, data []byte) {
		ts := time.Now().Format(time.RFC3339Nano)
		_, _ = f.WriteString(fmt.Sprintf("%s %d %s %s\n", ts, shard, topic, string(data)))
	}

	c.OnShardMessage(handleUpdate)

	return c, f, nil
}
