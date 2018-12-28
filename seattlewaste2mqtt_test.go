package main

import (
	"testing"
	"time"

	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
)

const knownGoodAddress = "2133 N 61ST ST"

func defaultSeattleWaste2Mqtt() *SeattleWaste2Mqtt {
	c := NewSeattleWaste2Mqtt(NewConfig(mqttExtCfg.NewMQTTConfig()), mqttExtDI.NewMQTTFuncWrapper())
	return c
}

func TestEncodeAddress(t *testing.T) {
	var tests = []struct {
		address       string
		encodeAddress string
	}{
		{knownGoodAddress, knownGoodAddress},
		{"12448 Fake Road Drive", ""},
	}

	c := defaultSeattleWaste2Mqtt()
	c.onConnect(c.client)

	for _, v := range tests {
		c.address = v.address
		c.encodedAddress = ""
		c.encodeAddress()
		if c.encodedAddress != v.encodeAddress {
			t.Errorf("Wrong encoded address. Actual: %s, Expected: %s", c.encodedAddress, v.encodeAddress)
		}
	}
}

func TestSeattleWaste2Mqtt(t *testing.T) {
	var tests = []struct {
		date string
	}{
		{"Mar 16th, 2017"},
		{"Aug 31st, 2017"},
		{"June 1st, 2017"},
	}

	c := defaultSeattleWaste2Mqtt()
	c.address = knownGoodAddress
	c.encodedAddress = knownGoodAddress

	layout := "2006-01-02"
	for _, v := range tests {
		now, _ := time.Parse(layout, v.date)
		collectionInfo, err := c.collectionLookup(now)
		if collectionInfo.Start == "" || err != nil {
			t.Errorf("Error looking up collection info")
		}
	}
}

func TestSeattleWaste2MqttLoop(t *testing.T) {
	c := defaultSeattleWaste2Mqtt()
	c.address = knownGoodAddress
	c.encodedAddress = knownGoodAddress
	c.loop(true)
}

func TestMqttRun(t *testing.T) {
	c := defaultSeattleWaste2Mqtt()
	if err := c.Run(); err != nil {
		t.Error("Something went wrong; expected to connect!")
	}

	c.client.Disconnect(0)
}

func TestMqttConnect(t *testing.T) {
	c := defaultSeattleWaste2Mqtt()
	c.onConnect(c.client)
}
