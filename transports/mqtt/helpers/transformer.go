package helpers

import (
	"fmt"
	"github.com/safecility/go/lib/stream"
)

type SimpleDaliPayloadAdjuster struct{}

func (i SimpleDaliPayloadAdjuster) AdjustPayload(message *stream.SimpleMessage) error {
	schema2Payload := append([]byte{0}, message.Payload...)
	message.Payload = schema2Payload
	return nil
}

type AppIdUidTransformer struct {
	AppID string
}

func (u AppIdUidTransformer) GetUID(deviceID string) string {
	return fmt.Sprintf("%s/%s", u.AppID, deviceID)
}
