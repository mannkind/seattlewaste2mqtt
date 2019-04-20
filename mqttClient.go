package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
	mqttExtHA "github.com/mannkind/paho.mqtt.golang.ext/ha"
)

type mqttClient struct {
	discovery       bool
	discoveryPrefix string
	discoveryName   string
	topicPrefix     string

	client mqtt.Client
}

func newMQTTClient(config *config, mqttFuncWrapper *mqttExtDI.MQTTFuncWrapper) *mqttClient {
	c := mqttClient{
		discovery:       config.MQTT.Discovery,
		discoveryPrefix: config.MQTT.DiscoveryPrefix,
		discoveryName:   config.MQTT.DiscoveryName,
		topicPrefix:     config.MQTT.TopicPrefix,
	}

	opts := mqttFuncWrapper.
		ClientOptsFunc().
		AddBroker(config.MQTT.Broker).
		SetClientID(config.MQTT.ClientID).
		SetOnConnectHandler(c.onConnect).
		SetConnectionLostHandler(c.onDisconnect).
		SetUsername(config.MQTT.Username).
		SetPassword(config.MQTT.Password).
		SetWill(c.availabilityTopic(), "offline", 0, true)

	c.client = mqttFuncWrapper.ClientFunc(opts)

	return &c
}

func (c *mqttClient) run() {
	log.Print("Connecting to MQTT")
	if token := c.client.Connect(); !token.Wait() || token.Error() != nil {
		log.Printf("Error connecting to MQTT: %s", token.Error())
		panic("Exiting...")
	}
}

func (c *mqttClient) onConnect(client mqtt.Client) {
	log.Print("Connected to MQTT")
	c.publish(c.availabilityTopic(), "online")
	c.publishDiscovery()
}

func (c *mqttClient) onDisconnect(client mqtt.Client, err error) {
	log.Printf("Disconnected from MQTT: %s.", err)
}

func (c *mqttClient) availabilityTopic() string {
	return fmt.Sprintf("%s/status", c.topicPrefix)
}

func (c *mqttClient) publishDiscovery() {
	if !c.discovery {
		return
	}

	obj := reflect.ValueOf(event{}.data)
	for i := 0; i < obj.NumField(); i++ {
		sensor := strings.ToLower(obj.Type().Field(i).Name)
		val := obj.Field(i)
		sensorType := ""

		switch val.Kind() {
		case reflect.Bool:
			sensorType = "binary_sensor"
		case reflect.String:
			sensorType = "sensor"
		}

		if sensorType == "" {
			continue
		}

		mqd := mqttExtHA.MQTTDiscovery{
			DiscoveryPrefix: c.discoveryPrefix,
			Component:       sensorType,
			NodeID:          c.discoveryName,
			ObjectID:        sensor,

			AvailabilityTopic: c.availabilityTopic(),
			Name:              fmt.Sprintf("%s %s", c.discoveryName, sensor),
			StateTopic:        fmt.Sprintf(sensorTopicTemplate, c.topicPrefix, sensor),
			UniqueID:          fmt.Sprintf("%s.%s", c.discoveryName, sensor),
		}

		mqd.PublishDiscovery(c.client)
	}
}

func (c *mqttClient) receive(e event) {
	info := e.data
	obj := reflect.ValueOf(info)
	for i := 0; i < obj.NumField(); i++ {
		sensor := strings.ToLower(obj.Type().Field(i).Name)
		val := obj.Field(i)

		topic := fmt.Sprintf(sensorTopicTemplate, c.topicPrefix, sensor)
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

		c.publish(topic, payload)
	}
}

func (c *mqttClient) publish(topic string, payload string) {
	retain := true
	if token := c.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}

	log.Print(fmt.Sprintf("Publishing - Topic: %s ; Payload: %s", topic, payload))
}
