package sections

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/helpers"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/messages"
)

func CheckOrCreateBigqueryTable(config *helpers.Config) (*bigquery.TableMetadata, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, config.ProjectName)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer func(client *bigquery.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close bigquery.Client")
		}
	}(client)

	tableRef := client.Dataset(config.BigQuery.Dataset).Table(config.BigQuery.Table)

	tableMetadata, err := tableRef.Metadata(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get table metadata")
	}
	if tableMetadata == nil {
		err = tableRef.Create(ctx, messages.GetBigqueryTableMetadata(config.BigQuery.Table))
		if err != nil {
			return nil, err
		}
		log.Info().Msg("Created bigquery table")
		tableMetadata, err = tableRef.Metadata(ctx)
	}

	return tableMetadata, err
}
