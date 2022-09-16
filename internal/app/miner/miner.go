package miner

import (
	"context"
	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/pubsub"
	"github.com/asaskevich/EventBus"
	"log"
	"os"
	"os/signal"
	"syscall"
	twitch2 "ttv-cli/internal/pkg/twitch"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
)

type Miner struct {
	UserId       string
	client       *twitch2.Client
	channels     []channelfollows.ChannelFollow
	pubsubClient *pubsub.Client
	eventBus     EventBus.Bus
}

func New(client *twitch2.Client, userId string) Miner {
	channels, err := channelfollows.Get(client)
	if err != nil {
		log.Fatalln("Unable to get followed channels: ", err)
	}

	return Miner{
		UserId:       userId,
		client:       client,
		pubsubClient: twitch.PubSub(),
		channels:     channels,
		eventBus:     EventBus.New(),
	}
}

func (m Miner) Start() {
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	if err := m.subscribeStreamStatus(ctx); err != nil {
		log.Println("could not subscribe to stream start events: ", err)
		goto exit
	}

	if err := m.subscribeMoments(ctx); err != nil {
		log.Println("could not subscribe to Moments events: ", err)
		goto exit
	}

	if err := m.subscribePoints(ctx); err != nil {
		log.Println("could not subscribe to points update events: ", err)
		goto exit
	}

	<-exitChannel

exit:
	cancel()
	log.Println("Exiting.")
}
