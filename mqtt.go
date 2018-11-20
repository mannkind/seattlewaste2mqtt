package main

import (
	"github.com/eclipse/paho.mqtt.golang"
)

type newMqttClientOptsFunc func() *mqtt.ClientOptions
type newMqttClientFunc func(*mqtt.ClientOptions) mqtt.Client
type mqttDiscovery struct {
	Name        string `json:"name"`
	StateTopic  string `json:"state_topic"`
	UniqueID    string `json:"unique_id,omitempty"`
	PayloadOn   string `json:"payload_on,omitempty"`
	PayloadOff  string `json:"payload_off,omitempty"`
	DeviceClass string `json:"device_class,omitempty"`
}

// MQTTFuncWrapper - Wraps the functions needed to create a new MQTT client.
type MQTTFuncWrapper struct {
	clientOptsFunc newMqttClientOptsFunc
	clientFunc     newMqttClientFunc
}

// NewMQTTFuncWrapper - Returns a fancy new wrapper for the mqtt creation functions.
func NewMQTTFuncWrapper() *MQTTFuncWrapper {
	return &MQTTFuncWrapper{
		clientOptsFunc: mqtt.NewClientOptions,
		clientFunc:     mqtt.NewClient,
	}
}
