module github.com/mannkind/seattlewaste2mqtt

go 1.13

require (
	github.com/caarlos0/env/v6 v6.0.0
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fatih/color v1.7.0 // indirect
	github.com/google/wire v0.4.0
	github.com/magefile/mage v1.9.0
	github.com/mannkind/seattlewaste v0.1.0
	github.com/mannkind/twomqtt v0.4.6
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/robfig/cron/v3 v3.0.0
	github.com/sirupsen/logrus v1.4.2
)

// local development
// replace github.com/mannkind/twomqtt => ../twomqtt
