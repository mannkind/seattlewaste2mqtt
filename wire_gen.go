// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/mannkind/twomqtt"
)

// Injectors from wire.go:

func initialize() *bridge {
	mainConfig := newConfig()
	mainMqttClientConfig := mainConfig.MQTTClientConfig
	mqttProxyConfig := mainMqttClientConfig.MQTTProxyConfig
	mqttProxy := twomqtt.NewMQTTProxy(mqttProxyConfig)
	mainMqttClient := newMQTTClient(mainMqttClientConfig, mqttProxy)
	mainServiceClientConfig := mainConfig.ServiceClientConfig
	mainServiceClient := newServiceClient(mainServiceClientConfig)
	mainBridge := newBridge(mainMqttClient, mainServiceClient)
	return mainBridge
}
