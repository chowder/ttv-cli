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
	"ttv-cli/internal/pkg/twitch/gql/operation/channelfollows"
)

type Miner struct {
	AuthToken    string
	UserId       string
	channels     []channelfollows.ChannelFollow
	pubsubClient *pubsub.Client
	eventBus     EventBus.Bus
}

func New(authToken string, userId string) Miner {
	channels, err := channelfollows.Get(authToken)
	if err != nil {
		log.Fatalln("Unable to get followed channels: ", err)
	}

	return Miner{
		AuthToken:    authToken,
		UserId:       userId,
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
