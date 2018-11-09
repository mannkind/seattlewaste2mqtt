# Seattle Waste MQTT

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/seattlewaste2mqtt/blob/master/LICENSE.md)
[![Travis CI](https://img.shields.io/travis/mannkind/seattlewaste2mqtt/master.svg?style=flat-square)](https://travis-ci.org/mannkind/seattlewaste2mqtt)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/seattlewaste2mqtt/master.svg)](http://codecov.io/github/mannkind/seattlewaste2mqtt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mannkind/seattlewaste2mqtt)](https://goreportcard.com/report/github.com/mannkind/seattlewaste2mqtt)

# Installation

## Via Docker
```
docker run -d --name="seattlewaste2mqtt" -v /etc/localtime:/etc/localtime:ro mannkind/seattlewaste2mqtt
```

## Via Make
```
git clone https://github.com/mannkind/seattlewaste2mqtt
cd seattlewaste2mqtt
make
SEATTLEWASTE_ADDRESS="2133 N 61ST ST" ./seattlewaste2mqtt 
```

# Configuration

Configuration happens via environmental variables

```
SEATTLEWASTE_ADDRESS - The address for which to lookup collections
SEATTLEWASTE_ALERTWITHIN - [OPTIONAL] The duration for which to alert, defaults to "24h"
SEATTLEWASTE_LOOKUPINTERVAL - [OPTIONAL] The duration for which to lookup collections, defaults to "8h"
MQTT_CLIENTID - [OPTIONAL] The clientId, defaults to "DefaultSeattleWaste2MQTTClientID"
MQTT_BROKER - [OPTIONAL] The MQTT broker, defaults to "tcp://mosquitto.org:1883"
MQTT_PUBTOPIC - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/seattle_waste"
MQTT_USERNAME - [OPTIONAL] The MQTT username, default to ""
MQTT_PASSWORD - [OPTIONAL] The MQTT password, default to ""
```
