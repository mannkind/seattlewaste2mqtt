//+build wireinject

package main

import (
	"github.com/google/wire"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
)

// Initialize - Compile-time DI
func Initialize() *SeattleWaste2Mqtt {
	wire.Build(mqttExtCfg.NewMQTTConfig, NewConfig, mqttExtDI.NewMQTTFuncWrapper, NewSeattleWaste2Mqtt)

	return &SeattleWaste2Mqtt{}
}
