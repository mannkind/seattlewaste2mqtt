package source

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/mannkind/seattlewaste2mqtt/shared"
	log "github.com/sirupsen/logrus"
)

// Opts is for package related settings
type Opts struct {
	shared.Opts
	AlertWithin    time.Duration `env:"SEATTLEWASTE_ALERTWITHIN" envDefault:"24h"`
	LookupInterval time.Duration `env:"SEATTLEWASTE_LOOKUPINTERVAL" envDefault:"8h"`
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

	return c
}
