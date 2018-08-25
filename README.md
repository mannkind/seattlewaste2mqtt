# Seattle Waste MQTT

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/seattle_waste_mqtt/blob/master/LICENSE.md)
[![Travis CI](https://img.shields.io/travis/mannkind/seattle_waste_mqtt/master.svg?style=flat-square)](https://travis-ci.org/mannkind/seattle_waste_mqtt)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/seattle_waste_mqtt/master.svg)](http://codecov.io/github/mannkind/seattle_waste_mqtt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mannkind/seattle_waste_mqtt)](https://goreportcard.com/report/github.com/mannkind/seattle_waste_mqtt)

# Installation

## Via Docker
```
docker run -d --name="seattle_waste_mqtt" -v /the/path/to/config_folder:/config -v /etc/localtime:/etc/localtime:ro mannkind/seattle_waste_mqtt
```

## Via Make
```
git clone https://github.com/mannkind/seattle_waste_mqtt
cd seattle_waste_mqtt
make
./bin/seattle_waste_mqtt -c */the/path/to/config_folder/config.yaml*
```

## Via Go
```
go get -u github.com/mannkind/seattle_waste_mqtt
go install github.com/mannkind/seattle_waste_mqtt
seattle_waste_mqtt -c */the/path/to/config_folder/config.yaml*
```

# Configuration

Configuration happens in the config.yaml file. A full example might look this:

```
settings:
    clientid: 'GoSeattleWasteMQTT'
    broker:   'tcp://mosquitto:1883'
    pubtopic: 'home/seattle_waste'

control:
    address: "Street Address w/o Seattle, WA and ZipCode"
    alertdaysbefore: 1
```
