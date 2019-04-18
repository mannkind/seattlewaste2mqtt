package main

import (
	"testing"

	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
)

const knownGoodAddress = "2133 N 61ST ST"

func defaultSeattleWaste2Mqtt() *SeattleWaste2Mqtt {
	c := NewSeattleWaste2Mqtt(NewConfig(mqttExtCfg.NewMQTTConfig()), mqttExtDI.NewMQTTFuncWrapper())
	return c
}

func TestSeattleWaste2MqttLoop(t *testing.T) {
	c := defaultSeattleWaste2Mqtt()
	c.address = knownGoodAddress
	c.loop(true)
}

func TestMqttConnect(t *testing.T) {
	c := defaultSeattleWaste2Mqtt()
	c.onConnect(c.client)
}
