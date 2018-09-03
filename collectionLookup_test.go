package main

import (
	"testing"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v2"
)

const knownGoodAddress = "2133 N 61ST ST"

var testClient = mqtt.NewClient(mqtt.NewClientOptions())

func defaultTestMQTT() *CollectionLookup {
	var testConfig = `
        settings:
          clientid: 'GoMySysBootloader'
          broker: "tcp://fake.mosquitto.org:1883"
          pubtopic: 'mysensors_tx'

        control:
            address: ''
    `

	myMqtt := CollectionLookup{}
	err := yaml.Unmarshal([]byte(testConfig), &myMqtt)
	if err != nil {
		panic(err)
	}
	return &myMqtt
}

func TestEncodeAddress(t *testing.T) {
	var tests = []struct {
		Address       string
		EncodeAddress string
	}{
		{knownGoodAddress, knownGoodAddress},
		{"12448 Fake Road Drive", ""},
	}

	myMQTT := defaultTestMQTT()
	myMQTT.onConnect(testClient)

	for _, v := range tests {
		myMQTT.Control.Address = v.Address
		myMQTT.Control.EncodedAddress = ""
		myMQTT.encodeAddress()
		if myMQTT.Control.EncodedAddress != v.EncodeAddress {
			t.Errorf("Wrong encoded address. Actual: %s, Expected: %s", myMQTT.Control.EncodedAddress, v.EncodeAddress)
		}
	}
}

func TestCollectionLookup(t *testing.T) {
	var tests = []struct {
		Date string
	}{
		{"Mar 16th, 2017"},
		{"Aug 31st, 2017"},
		{"June 1st, 2017"},
	}

	myMQTT := defaultTestMQTT()
	myMQTT.Control.Address = knownGoodAddress
	myMQTT.Control.EncodedAddress = knownGoodAddress

	layout := "2006-01-02"
	for _, v := range tests {
		now, err := time.Parse(layout, v.Date)
		collectionInfo, err := myMQTT.collectionLookup(now)
		if collectionInfo.Start == "" || err != nil {
			t.Errorf("Error looking up collection info")
		}
	}
}

func TestCollectionLookupLoop(t *testing.T) {
	myMQTT := defaultTestMQTT()
	myMQTT.Client = testClient
	myMQTT.Control.Address = knownGoodAddress
	myMQTT.Control.EncodedAddress = knownGoodAddress
	myMQTT.loop(true)
}

func TestMqttStart(t *testing.T) {
	myMQTT := defaultTestMQTT()
	if err := myMQTT.Start(); err == nil {
		t.Error("Something went wrong; expected a failure to connect!")
	}

	myMQTT.Stop()
}

func TestMqttConnect(t *testing.T) {
	myMQTT := defaultTestMQTT()
	myMQTT.onConnect(testClient)
}
