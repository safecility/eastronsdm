package helpers

import (
	"github.com/safecility/brokers/mqtt/messages"
	"github.com/safecility/go/lib/stream"
	"reflect"
	"testing"
)

func TestAppIdUidTransformer_GetUID(t *testing.T) {
	type fields struct {
		AppID string
	}
	type args struct {
		deviceID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AppIdUidTransformer{
				AppID: tt.fields.AppID,
			}
			if got := u.GetUID(tt.args.deviceID); got != tt.want {
				t.Errorf("GetUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegments2Scheme2Adjuster_AdjustPayload(t *testing.T) {
	type args struct {
		message *messages.LoraMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "",
			args: args{message: &messages.LoraMessage{
				SimpleMessage: stream.SimpleMessage{
					Payload: []byte{222, 222},
				},
			}},
			want:    []byte{2, 0, 222, 222},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.want, tt.args.message.Payload) {
				t.Errorf("AdjustPayload() = %v, want %v", tt.want, tt.args.message.Payload)
			}
		})
	}
}

func TestSimpleDaliPayloadAdjuster_AdjustPayload(t *testing.T) {
	type args struct {
		message *messages.LoraMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "",
			args: args{message: &messages.LoraMessage{
				SimpleMessage: stream.SimpleMessage{
					Payload: []byte{222, 222},
				},
			}},
			want:    []byte{2, 0, 222, 222},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
