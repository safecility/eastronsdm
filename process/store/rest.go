package store

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/safecility/go/lib"
	"net/url"
)

type DeviceClient struct {
	client *resty.Client
	server string
}

func CreateDeviceClient(serverAddress string) *DeviceClient {
	// Create a Resty Client
	client := resty.New()
	return &DeviceClient{
		client: client,
		server: serverAddress,
	}
}

func (dc *DeviceClient) GetDevice(uid string) (*lib.Device, error) {
	safeUID := url.PathEscape(uid)
	resp, err := dc.client.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s/device/%s", dc.server, safeUID))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%s", resp.Status())
	}
	pd := &lib.Device{}
	err = json.Unmarshal(resp.Body(), pd)
	return pd, err
}

func (dc *DeviceClient) Close() error {
	return nil
}
