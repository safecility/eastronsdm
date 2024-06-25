package messages

import (
	"cloud.google.com/go/bigquery"
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

func GetBigqueryTableMetadata(name string) *bigquery.TableMetadata {
	sampleSchema := bigquery.Schema{
		{Name: "DeviceUID", Type: bigquery.StringFieldType},
		{Name: "Time", Type: bigquery.TimestampFieldType},
		{Name: "ImportActiveEnergy", Type: bigquery.FloatFieldType},
		{Name: "ExportActiveEnergy", Type: bigquery.FloatFieldType},
		{Name: "ActivePower", Type: bigquery.FloatFieldType},
		{Name: "InstantaneousCurrent", Type: bigquery.FloatFieldType},
		{Name: "InstantaneousVoltage", Type: bigquery.FloatFieldType},
		{Name: "PowerFactor", Type: bigquery.FloatFieldType},
		{Name: "RelayState", Type: bigquery.FloatFieldType},
	}

	return &bigquery.TableMetadata{
		Name:   name,
		Schema: sampleSchema,
	}
}
