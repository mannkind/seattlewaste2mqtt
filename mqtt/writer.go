package mqtt

import (
	"reflect"
	"strings"

	"github.com/mannkind/seattlewaste2mqtt/shared"
	"github.com/mannkind/twomqtt"
)

// Writer is for writing a shared representation to MQTT
type Writer struct {
	*twomqtt.MQTT
	opts     Opts
	incoming <-chan shared.Representation
}

// NewWriter creates a new Writer for writing a shared representation to MQTT
func NewWriter(mqtt *twomqtt.MQTT, opts Opts, incoming <-chan shared.Representation) *Writer {
	c := Writer{
		MQTT:     mqtt,
		opts:     opts,
		incoming: incoming,
	}

	c.MQTT.
		SetDiscoveryHandler(c.discovery).
		SetReadIncomingChannelHandler(c.read).
		Initialize()

	return &c
}

func (c *Writer) discovery() []twomqtt.MQTTDiscovery {
	mqds := []twomqtt.MQTTDiscovery{}
	if !c.Discovery {
		return mqds
	}

	for _, deviceName := range c.opts.Addresses {
		obj := reflect.ValueOf(shared.Representation{})
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

			mqd := twomqtt.NewMQTTDiscovery(c.opts.MQTTOpts, deviceName, sensorName, sensorType)
			mqd.Device.Name = shared.Name
			mqd.Device.SWVersion = shared.Version

			mqds = append(mqds, *mqd)
		}
	}

	return mqds
}

func (c *Writer) read() {
	for info := range c.incoming {
		c.publish(info)
	}
}

func (c *Writer) publish(info shared.Representation) []twomqtt.MQTTMessage {
	published := []twomqtt.MQTTMessage{}

	name := c.opts.Addresses[info.Address]
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
