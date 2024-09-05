package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/helpers"
	"os"
	"time"
)

func main() {
	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to BigQuery")
	}
	defer func(client *bigquery.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close bigquery.Client")
		}
	}(client)

	bqc := gbigquery.NewBQTable(client)

	metaData := getTableMetadata(config.BigQuery.Table)
	t, err := bqc.CheckOrCreateBigqueryTable(&config.BigQuery, metaData)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create BigQuery table")
	}

	sClient, err := pubsub.NewSchemaClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create schema client")
	}
	defer func(sClient *pubsub.SchemaClient) {
		err := sClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close schema client")
		}
	}(sClient)

	schema, err := sClient.Schema(ctx, config.BigQuery.Schema.Name, pubsub.SchemaViewFull)
	if err != nil || schema == nil {
		schema, err = gbigquery.CreateProtoSchema(sClient, config.BigQuery.Schema.Name, config.BigQuery.Schema.FilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create schema")
		}
		log.Info().Str("revision", schema.RevisionID).Msg("Schema created")
	}
	if config.BigQuery.Schema.Revision != "" && schema.RevisionID != config.BigQuery.Schema.Revision {
		schema, err = gbigquery.UpdateProtoSchema(sClient, schema.Name, config.BigQuery.Schema.Revision, config.BigQuery.Schema.FilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("could not update schema")
		}
	}

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup pubsub")
	}

	bigqueryTopic := gpsClient.Topic(config.Pubsub.Topics.Bigquery)
	exists, err := bigqueryTopic.Exists(ctx)
	if !exists {
		bigqueryTopic, err = gbigquery.CreateBigqueryTopic(gpsClient, config.Pubsub.Topics.Bigquery, schema)
		if err != nil {
			log.Fatal().Str("sub", config.Pubsub.Subscriptions.BigQuery).Err(err).Msg("could not create bigquery topic")
		}
		log.Info().Msg("bigquery topic created")
	}
	bigQuerySubscription := gpsClient.Subscription(config.Pubsub.Subscriptions.BigQuery)
	exists, err = bigQuerySubscription.Exists(ctx)
	if !exists {
		err = gbigquery.CreateBigQuerySubscription(gpsClient, config.Pubsub.Subscriptions.BigQuery, t.FullID, bigqueryTopic)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create bigquery subscription")
		}
		log.Info().Msg("created bigquery subscription")
	}

	usageSubscription := gpsClient.Subscription(config.Pubsub.Subscriptions.Eastron)
	exists, err = usageSubscription.Exists(ctx)
	if !exists {
		usageTopic := gpsClient.Topic(config.Pubsub.Topics.Eastron)
		if exists, err = usageTopic.Exists(ctx); err != nil {
			log.Fatal().Err(err).Msg("could not check if eastron topic exists")
		}
		if !exists {
			usageTopic, err = gbigquery.CreateBigqueryTopic(gpsClient, config.Pubsub.Topics.Eastron, schema)
			if err != nil {
				log.Fatal().Err(err).Str("topic", config.Pubsub.Topics.Eastron).Msg("could not create eastron topic")
			}
			log.Info().Msg("created eastron topic")
		}

		r, err := time.ParseDuration("1h")
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse duration")
		}
		subConfig := stream.GetDefaultSubscriptionConfig(usageTopic, r)
		usageSubscription, err = gpsClient.CreateSubscription(ctx, config.Pubsub.Subscriptions.Eastron, subConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create subscription")
		}
		log.Info().Msg("created usage topic")
	}
	log.Info().Msg("finished pubsub setup")

}

func getTableMetadata(name string) *bigquery.TableMetadata {
	tableSchema := bigquery.Schema{
		{Name: "DeviceUID", Type: bigquery.StringFieldType},
		{Name: "Time", Type: bigquery.TimestampFieldType},
		{Name: "AccumulatedCurrent", Type: bigquery.FloatFieldType},
		{Name: "InstantaneousCurrent", Type: bigquery.FloatFieldType},
		{Name: "MaximumCurrent", Type: bigquery.FloatFieldType},
		{Name: "MinimumCurrent", Type: bigquery.FloatFieldType},
	}

	return &bigquery.TableMetadata{
		Name:   name,
		Schema: tableSchema,
	}
}
