package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/devices/eastronsdm/process/messages"
	"github.com/safecility/devices/eastronsdm/process/store"
	"github.com/safecility/go/lib"
	"github.com/safecility/go/lib/stream"
	"net/http"
	"os"
)

type EastronServer struct {
	cache        store.DeviceStore
	uplinks      *pubsub.Subscription
	eastronTopic *pubsub.Topic
	pipeAll      bool
}

func NewEastronServer(cache store.DeviceStore, u *pubsub.Subscription, t *pubsub.Topic, pipeAll bool) EastronServer {
	return EastronServer{uplinks: u, cache: cache, eastronTopic: t, pipeAll: pipeAll}
}

func (es *EastronServer) Start() {
	go es.receive()
	es.serverHttp()
}

func (es *EastronServer) receive() {

	err := es.uplinks.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		sm := &stream.SimpleMessage{}
		log.Debug().Str("data", fmt.Sprintf("%s", message.Data)).Msg("raw data")
		err := json.Unmarshal(message.Data, sm)
		message.Ack()
		if err != nil {
			log.Err(err).Msg("could not unmarshall data")
			return
		}
		eastronReading, err := messages.ReadEastronInfo(sm.Payload)
		if err != nil {
			log.Err(err).Msg("could not read payload")
			return
		}
		// the eastron uid is not actually very useful
		eastronReading.UID = sm.DeviceUID
		eastronReading.Time = sm.Time

		log.Debug().Str("eui", sm.DeviceUID).Msg("eastron message")
		var pd *lib.Device
		if es.cache != nil {
			pd, err = es.cache.GetDevice(sm.DeviceUID)
			if err != nil {
				log.Warn().Err(err).Str("uid", sm.DeviceUID).Msg("could not get device")
			}
			if pd == nil {
				log.Debug().Str("uid", sm.DeviceUID).Msg("device not found")
			}
			eastronReading.Device = pd
		}
		if eastronReading.Device == nil && !es.pipeAll {
			log.Debug().Str("device", sm.DeviceUID).Msg("no device in cache and pipeAll == false")
			return
		}

		topic, err := stream.PublishToTopic(eastronReading, es.eastronTopic)
		if err != nil {
			log.Err(err).Msg("could not publish usage to topic")
			return
		}
		log.Debug().Str("topic", *topic).Msg("published eastron sdm to topic")
	})
	if err != nil {
		log.Err(err).Msg("could not receive from sub")
		return
	}
}

func (es *EastronServer) serverHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
