package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/helpers"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/server"
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

	sub := gpsClient.Subscription(config.Pubsub.Subscriptions.Eastron)
	exists, err := sub.Exists(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failure on subscription exists")
	}
	if !exists {
		log.Fatal().Str("subscription", config.Pubsub.Subscriptions.Eastron).Msg("Subscription does not exist")
	}

	topic := gpsClient.Topic(config.Pubsub.Topics.Bigquery)
	exists, err = topic.Exists(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failure on topic exists")
	}
	if !exists {
		log.Fatal().Str("topic", config.Pubsub.Topics.Bigquery).Msg("Topic does not exist")
	}

	bigqueryServer := server.NewEastronServer(sub, topic, config.StoreAll)

	bigqueryServer.Start()
}
