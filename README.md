# seattlewaste2mqtt

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/seattlewaste2mqtt/blob/master/LICENSE.md)
[![Build Status](https://github.com/mannkind/seattlewaste2mqtt/workflows/Main%20Workflow/badge.svg)](https://github.com/mannkind/seattlewaste2mqtt/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/seattlewaste2mqtt/master.svg)](http://codecov.io/github/mannkind/seattlewaste2mqtt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mannkind/seattlewaste2mqtt)](https://goreportcard.com/report/github.com/mannkind/seattlewaste2mqtt)

## Installation

### Via Docker

```bash
docker run -d --name="seattlewaste2mqtt" -e "SEATTLEWASTE_ADDRESS=2133 N 61ST ST" -v /etc/localtime:/etc/localtime:ro mannkind/seattlewaste2mqtt
```

### Via Mage

```bash
git clone https://github.com/mannkind/seattlewaste2mqtt
cd seattlewaste2mqtt
mage
SEATTLEWASTE_ADDRESS="2133 N 61ST ST" ./seattlewaste2mqtt
```

## Configuration

Configuration happens via environmental variables

```bash
SEATTLEWASTE_ADDRESS        - The comma separated address:name pairs, defaults to ""
SEATTLEWASTE_ALERTWITHIN    - [OPTIONAL] The duration for which to alert, defaults to "24h"
SEATTLEWASTE_LOOKUPINTERVAL - [OPTIONAL] The duration for which to lookup collections, defaults to "8h"
MQTT_TOPICPREFIX            - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/seattle_waste"
MQTT_DISCOVERY              - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
MQTT_DISCOVERYPREFIX        - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
MQTT_DISCOVERYNAME          - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "seattle_waste"
MQTT_CLIENTID               - [OPTIONAL] The clientId, defaults to "DefaultSeattleWaste2MQTTClientWrapperID"
MQTT_BROKER                 - [OPTIONAL] The MQTT broker, defaults to "tcp://mosquitto.org:1883"
MQTT_USERNAME               - [OPTIONAL] The MQTT username, default to ""
MQTT_PASSWORD               - [OPTIONAL] The MQTT password, default to ""
```
