package main

import (
	"reflect"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

type mqttClient struct {
	mqttClientConfig
	*twomqtt.MQTTProxy
	stateUpdateChan stateChannel
}

func newMQTTClient(mqttClientCfg mqttClientConfig, client *twomqtt.MQTTProxy, stateUpdateChan stateChannel) *mqttClient {
	c := mqttClient{
		MQTTProxy:        client,
		mqttClientConfig: mqttClientCfg,
		stateUpdateChan:  stateUpdateChan,
	}

	c.Initialize(
		c.onConnect,
		c.onDisconnect,
	)

	c.LogSettings()

	return &c
}

func (c *mqttClient) run() {
	c.Run()
	go c.receive()
}

func (c *mqttClient) onConnect(client mqtt.Client) {
	log.Info("Finished connecting to MQTT")
	c.Publish(c.AvailabilityTopic(), "online")
	c.publishDiscovery()
}

func (c *mqttClient) onDisconnect(client mqtt.Client, err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error("Disconnected from MQTT")
}

func (c *mqttClient) publishDiscovery() {
	if !c.Discovery {
		return
	}

	log.Info("MQTT discovery publishing")

	for address, name := range c.Addresses {
		log.WithFields(log.Fields{
			"address": address,
		}).Debug("Iterating through addresses")

		obj := reflect.ValueOf(collection{})
		for i := 0; i < obj.NumField(); i++ {
			field := obj.Type().Field(i)
			sensor := strings.ToLower(field.Name)
			sensorOverride, sensorIgnored := twomqtt.MQTTOverride(field)
			sensorType, sensorTypeIgnored := twomqtt.MQTTDiscoveryOverride(field)

			// Skip any fields tagged as ignored
			if sensorIgnored || sensorTypeIgnored {
				continue
			}

			// Override sensor name
			if sensorOverride != "" {
				sensor = sensorOverride
			}

			mqd := c.NewMQTTDiscovery(name, sensor, sensorType)

			c.PublishDiscovery(mqd)
		}

		log.Debug("Finished iterating through addresses")
	}

	log.Info("Finished MQTT discovery publishing")
}

func (c *mqttClient) receive() {
	for info := range c.stateUpdateChan {
		c.receiveState(info)
	}
}

func (c *mqttClient) receiveState(info collection) {
	name := c.Addresses[info.Address]
	obj := reflect.ValueOf(info)

	log.WithFields(log.Fields{
		"info": info,
	}).Debug("Publishing received state")

	for i := 0; i < obj.NumField(); i++ {
		field := obj.Type().Field(i)
		val := obj.Field(i)
		sensor := strings.ToLower(field.Name)
		sensorOverride, sensorIgnored := twomqtt.MQTTOverride(field)
		_, sensorTypeIgnored := twomqtt.MQTTDiscoveryOverride(field)

		// Skip any fields tagged as ignored
		if sensorIgnored || sensorTypeIgnored {
			continue
		}

		// Override sensor name
		if sensorOverride != "" {
			sensor = sensorOverride
		}

		topic := c.StateTopic(name, sensor)
		payload := ""

		switch val.Kind() {
		case reflect.Bool:
			payload = "OFF"
			if val.Bool() {
				payload = "ON"
			}
		case reflect.String:
			payload = val.String()
		}

		if payload == "" {
			continue
		}

		c.Publish(topic, payload)
	}

	log.Debug("Finished publishing received state")
}
