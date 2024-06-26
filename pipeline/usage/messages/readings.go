package messages

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"time"
)

type EastronSdmReading struct {
	*lib.Device
	UID                  string
	ImportActiveEnergy   float64
	ExportActiveEnergy   float64
	ActivePower          float64
	InstantaneousCurrent float64
	InstantaneousVoltage float64
	PowerFactor          float64
	RelayState           float64
	Time                 time.Time
}

type MeterReading struct {
	*lib.Device
	ReadingKWH float64
	Time       time.Time
}

func (mc EastronSdmReading) Usage() (*MeterReading, error) {
	if mc.ImportActiveEnergy == 0 {
		log.Info().Str("reading", fmt.Sprintf("%+v", mc)).Msg("zero usage - check device is new")
	}
	kWh := mc.ImportActiveEnergy * mc.InstantaneousVoltage * mc.PowerFactor / 1000.0
	mr := &MeterReading{
		ReadingKWH: kWh,
		Time:       mc.Time,
	}
	if mc.Device == nil {
		log.Warn().Str("UID", mc.UID).Msg("device does not have device definitions")
		mr.Device = &lib.Device{
			DeviceUID: mc.UID,
		}
	} else {
		mr.Device = mc.Device
	}

	return mr, nil
}

type Alarms struct {
	t  bool
	tr bool
	r  bool
	rr bool
}

type Version struct {
	Ipso     string
	Hardware string
	Firmware string
}

type Current struct {
	Total float32
	Value float32
	Max   float32
	Min   float32
	Alarms
}
