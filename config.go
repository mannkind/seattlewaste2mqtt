package main

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
)

type config struct {
	MQTT           *mqttExtCfg.MQTTConfig
	Address        string        `env:"SEATTLEWASTE_ADDRESS,required"`
	AlertWithin    time.Duration `env:"SEATTLEWASTE_ALERTWITHIN"      envDefault:"24h"`
	LookupInterval time.Duration `env:"SEATTLEWASTE_LOOKUPINTERVAL"   envDefault:"8h"`
}

func newConfig(mqttCfg *mqttExtCfg.MQTTConfig) *config {
	c := config{}
	c.MQTT = mqttCfg

	if c.MQTT.ClientID == "" {
		c.MQTT.ClientID = "DefaultSeattleWaste2MQTTClientID"
	}

	if c.MQTT.DiscoveryName == "" {
		c.MQTT.DiscoveryName = "seattle_waste"
	}

	if c.MQTT.TopicPrefix == "" {
		c.MQTT.TopicPrefix = "home/seattle_waste"
	}

	if err := env.Parse(&c); err != nil {
		log.Printf("Error unmarshaling configuration: %s", err)
	}

	redactedPassword := ""
	if len(c.MQTT.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.Printf("Environmental Settings:")
	log.Printf("  * ClientID        : %s", c.MQTT.ClientID)
	log.Printf("  * Broker          : %s", c.MQTT.Broker)
	log.Printf("  * Username        : %s", c.MQTT.Username)
	log.Printf("  * Password        : %s", redactedPassword)
	log.Printf("  * Discovery       : %t", c.MQTT.Discovery)
	log.Printf("  * DiscoveryPrefix : %s", c.MQTT.DiscoveryPrefix)
	log.Printf("  * DiscoveryName   : %s", c.MQTT.DiscoveryName)
	log.Printf("  * TopicPrefix     : %s", c.MQTT.TopicPrefix)
	log.Printf("  * Address         : %s", c.Address)
	log.Printf("  * AlertWithin     : %s", c.AlertWithin)
	log.Printf("  * LookupInterval  : %s", c.LookupInterval)

	return &c
}
