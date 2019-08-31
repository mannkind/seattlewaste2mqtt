package main

import (
	"time"
)

type serviceClientConfig struct {
	globalClientConfig
	AlertWithin    time.Duration `env:"SEATTLEWASTE_ALERTWITHIN" envDefault:"24h"`
	LookupInterval time.Duration `env:"SEATTLEWASTE_LOOKUPINTERVAL" envDefault:"8h"`
}
