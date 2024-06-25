package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/iot/devices/eastronsdm/pipeline/messagestore/messages"
)

type DatastoreEastron struct {
	client *datastore.Client
}

func NewDatastoreEastron(client *datastore.Client) (*DatastoreEastron, error) {
	rd := &DatastoreEastron{client: client}
	return rd, nil
}

func (d *DatastoreEastron) AddEastronMessage(m *messages.EastronSdmReading) error {
	ctx := context.Background()
	k := datastore.IncompleteKey("EastronSdm", nil)
	k, err := d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}
	log.Debug().Str("uid", m.DeviceUID).Msg("putting new eastron eagle message")
	return nil
}
