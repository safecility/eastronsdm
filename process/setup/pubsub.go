package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/devices/eastronsdm/process/helpers"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/setup"
	"os"
	"time"
)

func main() {

	deployment, isSet := os.LookupEnv("Deployment")
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	ctx := context.Background()

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup pubsub")
	}

	eastronTopic := gpsClient.Topic(config.Topics.Eastron)
	exists, err := eastronTopic.Exists(ctx)
	if !exists {
		eastronTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Eastron)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create topic")
		}
		log.Info().Str("topic", eastronTopic.String()).Msg("created topic")
	}

	uSubscription := gpsClient.Subscription(config.Subscriptions.Uplinks)
	exists, err = uSubscription.Exists(ctx)
	if !exists {
		uTopic := gpsClient.Topic(config.Topics.Uplinks)
		exists, err = uTopic.Exists(ctx)

		if !exists {
			uTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Uplinks)
			if err != nil {
				log.Fatal().Err(err).Str("topic", config.Topics.Uplinks).Msg("setup could not create topic")
			}
			log.Info().Str("topic", uTopic.String()).Msg("created topic")
		}

		r, err := time.ParseDuration("1h")
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse duration")
		}
		subConfig := stream.GetDefaultSubscriptionConfig(uTopic, r)
		uSubscription, err = gpsClient.CreateSubscription(ctx, config.Subscriptions.Uplinks, subConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create subscription")
		}
		log.Info().Str("sub", uSubscription.String()).Msg("created subscription")

	}

	log.Info().Msg("setup complete")
}
