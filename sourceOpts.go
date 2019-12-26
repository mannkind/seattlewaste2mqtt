package main

import (
	"time"
)

type sourceOpts struct {
	globalOpts
	AlertWithin    time.Duration `env:"SEATTLEWASTE_ALERTWITHIN" envDefault:"24h"`
	LookupInterval time.Duration `env:"SEATTLEWASTE_LOOKUPINTERVAL" envDefault:"8h"`
}
