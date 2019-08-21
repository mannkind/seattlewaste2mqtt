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
	c.MQTT.Defaults("DefaultSeattleWaste2MQTTClientID", "seattle_waste", "home/seattle_waste")

	if err := env.Parse(&c); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	log.WithFields(log.Fields{
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
