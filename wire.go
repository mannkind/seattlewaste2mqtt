//+build wireinject

package main

import (
	"github.com/google/wire"
)

// Initialize - Compile-time DI
func Initialize() *SeattleWaste2Mqtt {
	wire.Build(NewConfig, NewMQTTFuncWrapper, NewSeattleWaste2Mqtt)

	return &SeattleWaste2Mqtt{}
}
