package main

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

const defaultDiscoveryName = "seattle_waste"
const defaultTopicPrefix = "home/seattle_waste"
const knownAddress = "135 NW 75th St"
const knownAddressName = "myhouse"
const knownDiscoveryName = "seattleWasteDiscoveryName"
const knownTopicPrefix = "home/seattleWasteMQTTTopicPrefix"
const knownOn = "ON"

var knownTypes = []string{
	"garbage", "recycling", "foodandyardwaste", "status",
}

func init() {
	log.SetLevel(log.PanicLevel)
}

func setEnvs(d, dn, tp, a string) {
	os.Setenv("MQTT_DISCOVERY", d)
	os.Setenv("MQTT_DISCOVERYNAME", dn)
	os.Setenv("MQTT_TOPICPREFIX", tp)
	os.Setenv("SEATTLEWASTE_ADDRESS", a)
}

func clearEnvs() {
	setEnvs("false", "", "", "")
}

func TestDiscovery(t *testing.T) {
	defer clearEnvs()

	for _, knownType := range knownTypes {
		var tests = []struct {
			Addresses       string
			DiscoveryName   string
			TopicPrefix     string
			ExpectedTopic   string
			ExpectedPayload string
		}{
			{
				knownAddress,
				defaultDiscoveryName,
				defaultTopicPrefix,
				"homeassistant/binary_sensor/" + defaultDiscoveryName + "/" + knownType + "/config",
				"{\"availability_topic\":\"" + defaultTopicPrefix + "/status\",\"device\":{\"identifiers\":[\"" + defaultTopicPrefix + "/status\"],\"manufacturer\":\"twomqtt\",\"name\":\"x2mqtt\",\"sw_version\":\"X.X.X\"},\"name\":\"" + defaultDiscoveryName + " " + knownType + "\",\"state_topic\":\"" + defaultTopicPrefix + "/" + knownType + "/state\",\"unique_id\":\"" + defaultDiscoveryName + "." + knownType + "\"}",
			},
			{
				knownAddress + ":" + knownAddressName,
				knownDiscoveryName,
				knownTopicPrefix,
				"homeassistant/binary_sensor/" + knownDiscoveryName + "/" + knownAddressName + "_" + knownType + "/config",
				"{\"availability_topic\":\"" + knownTopicPrefix + "/status\",\"device\":{\"identifiers\":[\"" + knownTopicPrefix + "/status\"],\"manufacturer\":\"twomqtt\",\"name\":\"x2mqtt\",\"sw_version\":\"X.X.X\"},\"name\":\"" + knownDiscoveryName + " " + knownAddressName + " " + knownType + "\",\"state_topic\":\"" + knownTopicPrefix + "/" + knownAddressName + "/" + knownType + "/state\",\"unique_id\":\"" + knownDiscoveryName + "." + knownAddressName + "." + knownType + "\"}",
			},
		}

		for _, v := range tests {
			setEnvs("true", v.DiscoveryName, v.TopicPrefix, v.Addresses)

			c := initialize()
			c.mqttClient.publishDiscovery()

			actualPayload := c.mqttClient.LastPublishedOnTopic(v.ExpectedTopic)
			if actualPayload != v.ExpectedPayload {
				t.Errorf("Actual:%s\nExpected:%s", actualPayload, v.ExpectedPayload)
			}
		}
	}
}

func TestReceieveState(t *testing.T) {
	defer clearEnvs()

	for _, knownType := range knownTypes {
		var tests = []struct {
			Addresses       string
			TopicPrefix     string
			ExpectedTopic   string
			ExpectedPayload string
		}{
			{
				knownAddress,
				defaultTopicPrefix,
				defaultTopicPrefix + "/" + knownType + "/state",
				knownOn,
			},
			{
				knownAddress + ":" + knownAddressName,
				knownTopicPrefix,
				knownTopicPrefix + "/" + knownAddressName + "/" + knownType + "/state",
				knownOn,
			},
		}

		obj := collection{
			Address:          knownAddress,
			Start:            "2019-08-01",
			Garbage:          true,
			Recycling:        true,
			FoodAndYardWaste: true,
			Status:           true,
		}

		for _, v := range tests {
			setEnvs("false", "", v.TopicPrefix, v.Addresses)
			c := initialize()
			c.mqttClient.receiveState(obj)

			actualPayload := c.mqttClient.LastPublishedOnTopic(v.ExpectedTopic)
			if actualPayload != v.ExpectedPayload {
				t.Errorf("Actual:%s\nExpected:%s", actualPayload, v.ExpectedPayload)
			}
		}
	}
}
