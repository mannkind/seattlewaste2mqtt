package mqtt

import (
	"github.com/caarlos0/env/v6"
	"github.com/mannkind/seattlewaste2mqtt/shared"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

// Opts is for package related settings
type Opts struct {
	shared.Opts
	MQTTOpts twomqtt.MQTTOpts
}

// NewOpts creates a Opts based on environment variables
func NewOpts(opts shared.Opts) Opts {
	c := Opts{
		Opts: opts,
	}

	if err := env.Parse(&c); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	// Defaults
	if c.MQTTOpts.DiscoveryName == "" {
		c.MQTTOpts.DiscoveryName = "seattle_waste"
	}

	if c.MQTTOpts.TopicPrefix == "" {
		c.MQTTOpts.TopicPrefix = "home/seattle_waste"
	}

	return c
}
