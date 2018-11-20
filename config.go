package main

import (
	"log"
	"time"

	"github.com/caarlos0/env"
)

// Config - Structured configuration for the application.
type Config struct {
	ClientID        string        `env:"MQTT_CLIENTID" envDefault:"DefaultSeattleWaste2MQTTClientID"`
	Broker          string        `env:"MQTT_BROKER" envDefault:"tcp://mosquitto.org:1883"`
	Username        string        `env:"MQTT_USERNAME"`
	Password        string        `env:"MQTT_PASSWORD"`
	Discovery       bool          `env:"SEATTLEWASTE_DISCOVERY" envDefault:"false"`
	DiscoveryPrefix string        `env:"SEATTLEWASTE_DISCOVERYPREFIX" envDefault:"homeassistant"`
	DiscoveryName   string        `env:"SEATTLEWASTE_DISCOVERYNAME" envDefault:"seattle_waste"`
	PubTopic        string        `env:"SEATTLEWASTE_PUBTOPIC" envDefault:"home/seattle_waste"`
	Address         string        `env:"SEATTLEWASTE_ADDRESS,required"`
	AlertWithin     time.Duration `env:"SEATTLEWASTE_ALERTWITHIN" envDefault:"24h"`
	LookupInterval  time.Duration `env:"SEATTLEWASTE_LOOKUPINTERVAL" envDefault:"8h"`
}

// NewConfig - Returns a new Config object with configuration from ENV variables.
func NewConfig() *Config {
	c := Config{}

	if err := env.Parse(&c); err != nil {
		log.Printf("Error unmarshaling configuration: %s", err)
	}

	redactedPassword := ""
	if len(c.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.Printf("Environmental Settings:")
	log.Printf("  * ClientID        : %s", c.ClientID)
	log.Printf("  * Broker          : %s", c.Broker)
	log.Printf("  * Username        : %s", c.Username)
	log.Printf("  * Password        : %s", redactedPassword)
	log.Printf("  * Discovery       : %t", c.Discovery)
	log.Printf("  * DiscoveryPrefix : %s", c.DiscoveryPrefix)
	log.Printf("  * DiscoveryName   : %s", c.DiscoveryName)
	log.Printf("  * PubTopic        : %s", c.PubTopic)
	log.Printf("  * Address         : %s", c.Address)
	log.Printf("  * AlertWithin     : %s", c.AlertWithin)
	log.Printf("  * LookupInterval  : %s", c.LookupInterval)

	return &c
}
