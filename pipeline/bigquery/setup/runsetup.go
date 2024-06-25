package main

import (
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/helpers"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/setup/sections"
	"os"
)

func main() {
	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	tmd, err := sections.CheckOrCreateBigqueryTable(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating bigquery table")
	}
	sections.SetupPubsub(config, tmd)

}
