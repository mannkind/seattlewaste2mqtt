# Seattle Waste MQTT

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/seattlewaste2mqtt/blob/master/LICENSE.md)
[![Travis CI](https://img.shields.io/travis/mannkind/seattlewaste2mqtt/master.svg?style=flat-square)](https://travis-ci.org/mannkind/seattlewaste2mqtt)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/seattlewaste2mqtt/master.svg)](http://codecov.io/github/mannkind/seattlewaste2mqtt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mannkind/seattlewaste2mqtt)](https://goreportcard.com/report/github.com/mannkind/seattlewaste2mqtt)

# Installation

## Via Docker
```
docker run -d --name="seattlewaste2mqtt" -v /the/path/to/config_folder:/config -v /etc/localtime:/etc/localtime:ro mannkind/seattlewaste2mqtt
```

## Via Make
```
git clone https://github.com/mannkind/seattlewaste2mqtt
cd seattlewaste2mqtt
make
./bin/seattlewaste2mqtt -c */the/path/to/config_folder/config.yaml*
```

# Configuration

Configuration happens in the config.yaml file. A full example might look this:

```
settings:
    clientid: 'GoSeattleWasteMQTT'
    broker:   'tcp://mosquitto:1883'
    pubtopic: 'home/seattle_waste'

control:
    address: '<street address without city/state/zip>'
    alertwithin: '24h'
    lookupinterval: '8h'
```
