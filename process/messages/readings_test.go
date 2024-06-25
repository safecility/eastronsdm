package messages

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"testing"
)

func isClose(test, expected, difference float64) bool {
	return test-expected < difference
}

func TestBytesToFloat32(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToFloat32(tt.args.b); got != tt.want {
				t.Errorf("BytesToFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadEastronInfo(t *testing.T) {

	good1, err := base64.StdEncoding.DecodeString("AVaXWQEYQxPsiwAAAABB8doMQ2igbz90hmY+azt21jU=")
	good2, err := base64.StdEncoding.DecodeString("AVaR5AEYRIwCjwAAAABDcCnRQ1qG6z9/FtM/la6/kX8=")
	new1, err := base64.StdEncoding.DecodeString("AVaUegEYRYt0ZgAAAABEbV3PQ1n3pD9+jo9AjA2G4qc=")
	misconfigured, err := base64.StdEncoding.DecodeString("AVaXHQEYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJSc=")
	if err != nil {
		log.Err(err).Msg("Error decoding recent")
		return
	}

	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *EastronReading
		wantErr bool
	}{
		// TODO: Add more better test cases.
		{name: "good1", args: args{payload: good1}, want: &EastronReading{
			ImportActiveEnergy: 148, ActivePower: 30, InstantaneousVoltage: 232, PowerFactor: 0.955,
		}},
		{name: "good2", args: args{payload: good2}, want: &EastronReading{
			ImportActiveEnergy: 1120, ActivePower: 240, InstantaneousVoltage: 218, PowerFactor: 0.995,
		}},
		{name: "new1", args: args{payload: new1}, want: &EastronReading{
			ImportActiveEnergy: 1120, ActivePower: 240, InstantaneousVoltage: 218, PowerFactor: 0.995,
		}},
		{name: "misconfigured", args: args{payload: misconfigured}, want: &EastronReading{
			ImportActiveEnergy: 0, ActivePower: 0, InstantaneousVoltage: 0, PowerFactor: 0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadEastronInfo(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEastronInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tErr := !isClose(got.ActivePower, tt.want.ActivePower, 1) && !isClose(got.PowerFactor, tt.want.PowerFactor, 0.01)
			if tErr {
				t.Errorf("ReadEastronInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
