package cmd

import (
	"context"
	"fmt"
	"github.com/Adeithe/go-twitch"
	"github.com/spf13/cobra"
	"log"
	"ttv-cli/internal/app/miner/moments"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql/query/users"
)

var minerCmd = &cobra.Command{
	Use:   "miner",
	Short: "Mines channel points, moments, and drops",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(context.Background(), args)
	},
}

func run(ctx context.Context, names []string) {
	cfg := config.CreateOrRead()

	p := twitch.PubSub()

	p.OnShardConnect(func(shard int) {
		fmt.Printf("Shard #%d connected!\n", shard)
	})

	streamerByIds := make(map[string]string, 0)
	us, err := users.GetUsers(names)
	if err != nil {
		log.Fatalf("Could not fetch channel information for users - %s\n", err)
	}

	for i, u := range us {
		if len(u.Id) == 0 {
			log.Printf("Could not find user with name '%s'\n", names[i])
			continue
		}
		streamerByIds[u.Id] = u.DisplayName
	}

	err = moments.MineMoments(p, streamerByIds, cfg.AuthToken)
	if err != nil {
		log.Fatalf("Could not subscribe to Moments - %s\n", err)
	}

	defer p.Close()

	fmt.Printf("Started listening to %d topics on %d shards\n", p.GetNumTopics(), p.GetNumShards())

	<-ctx.Done()
}

func init() {
	rootCmd.AddCommand(minerCmd)
}
