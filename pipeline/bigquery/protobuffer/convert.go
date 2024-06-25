package protobuffer

import (
	"github.com/safecility/iot/devices/eastronsdm/pipeline/bigquery/messages"
	"time"
)

func CreateProtobufMessage(r *messages.EastronSdmReading) *EastronSdmBq {
	bq := &EastronSdmBq{
		DeviceUID:            r.UID,
		Time:                 r.Time.Format(time.RFC3339),
		ActivePower:          r.ActivePower,
		ImportActiveEnergy:   r.ImportActiveEnergy,
		ExportActiveEnergy:   r.ExportActiveEnergy,
		InstantaneousCurrent: r.InstantaneousCurrent,
		InstantaneousVoltage: r.InstantaneousVoltage,
		PowerFactor:          r.PowerFactor,
	}
	if r.Device != nil {
		bq.DeviceUID = r.DeviceUID
	}
	return bq
}
