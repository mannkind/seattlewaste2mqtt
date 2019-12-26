package main

import (
	"reflect"
	"strings"

	"github.com/mannkind/twomqtt"
)

type sink struct {
	*twomqtt.MQTT
	config   sinkOpts
	incoming <-chan sourceRep
}

func newSink(mqtt *twomqtt.MQTT, config sinkOpts, incoming <-chan sourceRep) *sink {
	c := sink{
		MQTT:     mqtt,
		config:   config,
		incoming: incoming,
	}

	c.MQTT.
		SetDiscoveryHandler(c.discovery).
		SetReadIncomingChannelHandler(c.read).
		Initialize()

	return &c
}

func (c *sink) run() {
	c.Run()
}

func (c *sink) discovery() []twomqtt.MQTTDiscovery {
	mqds := []twomqtt.MQTTDiscovery{}
	if !c.Discovery {
		return mqds
	}

	for _, deviceName := range c.config.Addresses {
		obj := reflect.ValueOf(sourceRep{})
		for i := 0; i < obj.NumField(); i++ {
			field := obj.Type().Field(i)
			sensorName := strings.ToLower(field.Name)
			sensorOverride, sensorIgnored := twomqtt.MQTTOverride(field)
			sensorType, sensorTypeIgnored := twomqtt.MQTTDiscoveryOverride(field)

			// Skip any fields tagged as ignored
			if sensorIgnored || sensorTypeIgnored {
				continue
			}

			// Override sensor name
			if sensorOverride != "" {
				sensorName = sensorOverride
			}

			mqd := twomqtt.NewMQTTDiscovery(c.config.MQTTOpts, deviceName, sensorName, sensorType)
			mqd.Device.Name = Name
			mqd.Device.SWVersion = Version

			mqds = append(mqds, *mqd)
		}
	}

	return mqds
}

func (c *sink) read() {
	for info := range c.incoming {
		c.publish(info)
	}
}

func (c *sink) publish(info sourceRep) []twomqtt.MQTTMessage {
	published := []twomqtt.MQTTMessage{}

	name := c.config.Addresses[info.Address]
	obj := reflect.ValueOf(info)

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

		msg := c.Publish(topic, payload)
		published = append(published, msg)
	}

	return published
}
