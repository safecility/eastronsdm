package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/usage/helpers"
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

	usageTopic := gpsClient.Topic(config.Topics.Usage)
	exists, err := usageTopic.Exists(ctx)
	if !exists {
		usageTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Usage)
		if err != nil {
			log.Fatal().Err(err).Str("topic", config.Topics.Usage).Msg("setup could not create topic")
		}
		log.Info().Str("topic", usageTopic.String()).Msg("created topic")
	}

	eastronSubscription := gpsClient.Subscription(config.Subscriptions.Eastron)
	exists, err = eastronSubscription.Exists(ctx)
	if !exists {
		eastronTopic := gpsClient.Topic(config.Topics.Eastron)
		exists, err = eastronTopic.Exists(ctx)
		if !exists {
			eastronTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Eastron)
			if err != nil {
				log.Fatal().Err(err).Str("topic", config.Topics.Eastron).Msg("setup could not create topic")
			}
			log.Info().Str("topic", eastronTopic.String()).Msg("created topic")
		}

		r, err := time.ParseDuration("1h")
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse duration")
		}
		subConfig := stream.GetDefaultSubscriptionConfig(eastronTopic, r)
		eastronSubscription, err = gpsClient.CreateSubscription(ctx, config.Subscriptions.Eastron, subConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create subscription")
		}
		log.Info().Str("topic", eastronSubscription.String()).Msg("created subscription")
	}

	log.Info().Msg("setup complete")
}
