package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/messagestore/helpers"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/messagestore/server"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/messagestore/store"
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
	defer func(gpsClient *pubsub.Client) {
		err := gpsClient.Close()
		if err != nil {
			log.Err(err).Msg("Error closing pubsub client")
		}
	}(gpsClient)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if gpsClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return // this is here so golang doesn't complain about gpsClient being nil
	}

	hotdropSubscription := gpsClient.Subscription(config.Subscriptions.Eastron)
	exists, err := hotdropSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Eastron).Msg("no eastron subscription")
	}

	dsClient, err := datastore.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start service")
	}
	d, err := store.NewDatastoreEastron(dsClient)

	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore hotdrop")
	}

	hotDropServer := server.NewEastronServer(d, hotdropSubscription, config.StoreAll)
	hotDropServer.Start()
}
