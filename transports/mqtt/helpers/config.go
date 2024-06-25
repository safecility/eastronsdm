package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"os"
)

const (
	OSDeploymentKey = "LORA_DEPLOYMENT"
)

type MqttConfig struct {
	AppID    string `json:"appID"`
	Username string `json:"username"`
	Address  string `json:"address"`
	Downlink bool   `json:"downlink"`
	Location bool   `json:"location"`
	Signal   bool   `json:"signal"`
}

type Config struct {
	ProjectName string `json:"projectName"`
	Mqtt        MqttConfig
	Secret      setup.Secret `json:"secret"`
	Topics      struct {
		Joins            string `json:"joins"`
		Uplinks          string `json:"uplinks"`
		Downlinks        string `json:"downlinks"`
		DownlinkReceipts string `json:"downlinkReceipts"`
		Location         string `json:"location"`
		Signal           string `json:"signal"`
	} `json:"topics"`
	Subscriptions struct {
		Downlinks string `json:"downlinks"`
	} `json:"subscriptions"`
}

// GetConfig for ttn the Username has the form: username = fmt.Sprintf("%s@ttn", p.AppID)
func GetConfig(deployment string) *Config {
	fileName := fmt.Sprintf("%s-config.json", deployment)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not find config file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Err(err).Msg("could not close config during defer")
		}
	}(file)
	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal().Err(err).Str("filename", fileName).Msg("could not decode pubsub config")
	}
	return config
}
