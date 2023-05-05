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
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
)

type Miner struct {
	UserId       string
	client       config.Config
	channels     []channelfollows.ChannelFollow
	pubsubClient *pubsub.Client
	eventBus     EventBus.Bus
}

func New(config config.Config) Miner {
	channels, err := channelfollows.Get(config)
	if err != nil {
		log.Fatalln("Unable to get followed channels: ", err)
	}

	return Miner{
		UserId:       "",
		client:       config,
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
