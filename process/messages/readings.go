package messages

import (
	"encoding/binary"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"math"
	"time"
)

type EastronReading struct {
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

func ReadEastronInfo(payload []byte) (*EastronReading, error) {
	var by []byte

	dpi := &EastronReading{}

	if len(payload) < 4 {
		return nil, fmt.Errorf("message too short %d", len(payload))
	} else {
		by = payload[:4]
	}

	beI := binary.BigEndian.Uint32(by)
	dpi.UID = fmt.Sprintf("%d", beI)

	if len(payload) > 9 {
		by = payload[6:10]
	} else {
		return dpi, fmt.Errorf("expected len > 9, %d", len(payload))
	}
	dpi.ImportActiveEnergy = float64(BytesToFloat32(by))

	if len(payload) > 13 {
		by = payload[10:14]
	} else {
		return dpi, fmt.Errorf("expected len > 13, %d", len(payload))
	}
	dpi.ExportActiveEnergy = float64(BytesToFloat32(by))

	if len(payload) > 17 {
		by = payload[14:18]
	} else {
		return dpi, fmt.Errorf("expected len > 17, %d", len(payload))
	}
	dpi.ActivePower = float64(BytesToFloat32(by))

	if len(payload) > 21 {
		by = payload[18:22]
	} else {
		return dpi, fmt.Errorf("expected len > 21, %d", len(payload))
	}
	dpi.InstantaneousVoltage = float64(BytesToFloat32(by))

	if len(payload) > 25 {
		by = payload[22:26]
	} else {
		return dpi, fmt.Errorf("expected len > 25, %d", len(payload))
	}
	dpi.PowerFactor = float64(BytesToFloat32(by))

	if len(payload) > 29 {
		by = payload[26:30]
	} else {
		return dpi, fmt.Errorf("expected len > 29, %d", len(payload))
	}
	dpi.InstantaneousCurrent = float64(BytesToFloat32(by))

	log.Debug().Interface("data", dpi).Msg("read eastron data")
	return dpi, nil
}

func BytesToFloat32(b []byte) float32 {
	beI := binary.BigEndian.Uint32(b)

	return math.Float32frombits(beI)
}
