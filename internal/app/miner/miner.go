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
	channels     []channelfollows.ChannelFollow
	pubsubClient *pubsub.Client
	eventBus     EventBus.Bus
}

func New(authToken string) Miner {
	channels, err := channelfollows.Get(authToken)
	if err != nil {
		log.Fatalln("Unable to get followed channels: ", err)
	}

	return Miner{
		AuthToken:    authToken,
		pubsubClient: twitch.PubSub(),
		channels:     channels,
		eventBus:     EventBus.New(),
	}
}

func (m Miner) Start() {
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	if err := m.subscribeStreamStart(ctx); err != nil {
		log.Println("could not subscribe to stream start events: ", err)
		goto exit
	}

	if err := m.listenStreamStart(ctx); err != nil {
		log.Println("could not listen to stream start events: ", err)
		goto exit
	}

	if err := m.subscribeMoments(ctx); err != nil {
		log.Println("could not subscribe to Moments events: ", err)
		goto exit
	}

	if err := m.listenMoments(ctx); err != nil {
		log.Println("could not listen to Moments events: ", err)
		goto exit
	}

	<-exitChannel

exit:
	cancel()
	log.Println("Exiting.")
}
