package helpers

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/iot/devices/eastronsdm/process/store"
)

func GetStore(config *Config) (store.DeviceStore, error) {
	if config.Store.Rest != nil {
		return store.CreateDeviceClient(config.Store.Rest.Address()), nil
	}
	return getSql(config)
}

func getSql(config *Config) (*store.DeviceSql, error) {
	log.Warn().Msg("sql is no longer supported - we'll switch to firestore in here soon and prefer rest for live deployment")

	ctx := context.Background()

	secretsClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create secrets client")
	}
	defer func(secretsClient *secretmanager.Client) {
		err := secretsClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close secrets client")
		}
	}(secretsClient)

	sqlSecret := setup.GetNewSecrets(config.ProjectName, secretsClient)
	password, err := sqlSecret.GetSecret(config.Store.Sql.Secret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get secret")
	}
	config.Store.Sql.Config.Password = string(password)

	s, err := setup.NewSafecilitySql(config.Store.Sql.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility sql")
	}
	c, err := store.NewDeviceSql(s)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility device sql")
	}

	return c, err
}
