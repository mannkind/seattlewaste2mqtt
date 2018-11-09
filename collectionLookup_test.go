package main

import (
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/eclipse/paho.mqtt.golang"
)

const knownGoodAddress = "2133 N 61ST ST"

var testClient = mqtt.NewClient(mqtt.NewClientOptions())

func defaultTestMQTT() *collectionLookup {
	myMqtt := collectionLookup{}
	env.Parse(&myMqtt)
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
		myMQTT.Address = v.Address
		myMQTT.encodedAddress = ""
		myMQTT.encodeAddress()
		if myMQTT.encodedAddress != v.EncodeAddress {
			t.Errorf("Wrong encoded address. Actual: %s, Expected: %s", myMQTT.encodedAddress, v.EncodeAddress)
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
	myMQTT.Address = knownGoodAddress
	myMQTT.encodedAddress = knownGoodAddress

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
	myMQTT.Address = knownGoodAddress
	myMQTT.client = testClient
	myMQTT.encodedAddress = knownGoodAddress
	myMQTT.loop(true)
}

func TestMqttStart(t *testing.T) {
	myMQTT := defaultTestMQTT()
	if err := myMQTT.start(); err != nil {
		t.Error("Something went wrong; expected to connect!")
	}

	myMQTT.stop()
}

func TestMqttConnect(t *testing.T) {
	myMQTT := defaultTestMQTT()
	myMQTT.onConnect(testClient)
}
