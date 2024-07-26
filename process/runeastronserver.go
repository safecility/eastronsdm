package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/eastronsdm/process/helpers"
	"github.com/safecility/iot/devices/eastronsdm/process/server"
	"os"
)

func main() {

	ctx := context.Background()

	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if gpsClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return
	}

	uplinksSubscription := gpsClient.Subscription(config.Subscriptions.Uplinks)
	exists, err := uplinksSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Uplinks).Msg("no uplinks subscription")
	}

	eastronTopic := gpsClient.Topic(config.Topics.Eastron)
	exists, err = eastronTopic.Exists(ctx)
	if !exists {
		log.Fatal().Str("topic", config.Topics.Eastron).Msg("no hotdrop topic")
	}
	if err != nil {
		log.Fatal().Err(err).Str("topic", config.Topics.Eastron).Msg("could not get topic")
	}
	defer eastronTopic.Stop()

	s, err := helpers.GetStore(config)

	if err != nil {
		log.Fatal().Err(err).Msg("could not get store")
	}

	eagleServer := server.NewEastronServer(s, uplinksSubscription, eastronTopic, config.PipeAll)
	eagleServer.Start()

}
