package main

import (
	"time"

	"github.com/caarlos0/env"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	log "github.com/sirupsen/logrus"
)

type config struct {
	MQTT           *mqttExtCfg.MQTTConfig
	Address        string        `env:"SEATTLEWASTE_ADDRESS,required"`
	AlertWithin    time.Duration `env:"SEATTLEWASTE_ALERTWITHIN"      envDefault:"24h"`
	LookupInterval time.Duration `env:"SEATTLEWASTE_LOOKUPINTERVAL"   envDefault:"8h"`
	DebugLogLevel  bool          `env:"SEATTLEWASTE_DEBUG" envDefault:"false"`
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
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	redactedPassword := ""
	if len(c.MQTT.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.WithFields(log.Fields{
		"MQTT.ClientID":               c.MQTT.ClientID,
		"MQTT.Broker":                 c.MQTT.Broker,
		"MQTT.Username":               c.MQTT.Username,
		"MQTT.Password":               redactedPassword,
		"MQTT.Discovery":              c.MQTT.Discovery,
		"MQTT.DiscoveryPrefix":        c.MQTT.DiscoveryPrefix,
		"MQTT.DiscoveryName":          c.MQTT.DiscoveryName,
		"MQTT.TopicPrefix":            c.MQTT.TopicPrefix,
		"SeattleWaste.Address":        c.Address,
		"SeattleWaste.AlertWithin":    c.AlertWithin,
		"SeattleWaste.LookupInterval": c.LookupInterval,
		"SeattleWaste.DebugLogLevel":  c.DebugLogLevel,
	}).Info("Environmental Settings")

	if c.DebugLogLevel {
		log.SetLevel(log.DebugLevel)
		log.Debug("Enabling the debug log level")
	}

	return &c
}
