//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mannkind/seattlewaste2mqtt/mqtt"
	"github.com/mannkind/seattlewaste2mqtt/shared"
	"github.com/mannkind/seattlewaste2mqtt/source"
	"github.com/mannkind/twomqtt"
)

func initialize() *app {
	wire.Build(
		newApp,
		shared.NewOpts,
		shared.NewRepresentationChannel,
		shared.NewRepresentationChannelIncoming,
		shared.NewRepresentationChannelOutgoing,
		mqtt.NewOpts,
		mqtt.NewWriter,
		source.NewOpts,
		source.NewService,
		source.NewReader,
		wire.FieldsOf(new(mqtt.Opts), "MQTTOpts"),
		twomqtt.NewMQTT,
	)

	return &app{}
}
