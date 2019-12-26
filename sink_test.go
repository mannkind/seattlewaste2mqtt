package main

import (
	"os"
	"testing"

	"github.com/mannkind/twomqtt"
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
			Addresses          string
			DiscoveryName      string
			TopicPrefix        string
			ExpectedName       string
			ExpectedStateTopic string
			ExpectedUniqueID   string
		}{
			{
				knownAddress,
				defaultDiscoveryName,
				defaultTopicPrefix,
				defaultDiscoveryName + " " + knownType,
				defaultTopicPrefix + "/" + knownType + "/state",
				defaultDiscoveryName + "." + knownType,
			},
			{
				knownAddress + ":" + knownAddressName,
				knownDiscoveryName,
				knownTopicPrefix,
				knownDiscoveryName + " " + knownAddressName + " " + knownType,
				knownTopicPrefix + "/" + knownAddressName + "/" + knownType + "/state",
				knownDiscoveryName + "." + knownAddressName + "." + knownType,
			},
		}

		for _, v := range tests {
			setEnvs("true", v.DiscoveryName, v.TopicPrefix, v.Addresses)

			c := initialize()
			mqds := c.sink.discovery()

			mqd := twomqtt.MQTTDiscovery{}
			for _, tmqd := range mqds {
				if tmqd.Name == v.ExpectedName {
					mqd = tmqd
					break
				}
			}

			if mqd.Name != v.ExpectedName {
				t.Errorf("discovery Name does not match; %s vs %s", mqd.Name, v.ExpectedName)
			}
			if mqd.StateTopic != v.ExpectedStateTopic {
				t.Errorf("discovery StateTopic does not match; %s vs %s", mqd.StateTopic, v.ExpectedStateTopic)
			}
			if mqd.UniqueID != v.ExpectedUniqueID {
				t.Errorf("discovery UniqueID does not match; %s vs %s", mqd.UniqueID, v.ExpectedUniqueID)
			}
		}
	}
}

func TestPublish(t *testing.T) {
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

		obj := sourceRep{
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

			allPublished := c.sink.publish(obj)

			matching := twomqtt.MQTTMessage{}
			for _, state := range allPublished {
				if state.Topic == v.ExpectedTopic {
					matching = state
					break
				}
			}

			if matching.Payload != v.ExpectedPayload {
				t.Errorf("Actual:%s\nExpected:%s", matching.Payload, v.ExpectedPayload)
			}
		}
	}
}
