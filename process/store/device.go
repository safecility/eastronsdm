package store

import (
	"github.com/safecility/go/lib"
)

type DeviceStore interface {
	GetDevice(uid string) (*lib.Device, error)
}
