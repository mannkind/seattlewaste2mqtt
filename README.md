# seattlewaste2mqtt

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/seattlewaste2mqtt/blob/master/LICENSE.md)
[![Build Status](https://github.com/mannkind/seattlewaste2mqtt/workflows/Main%20Workflow/badge.svg)](https://github.com/mannkind/seattlewaste2mqtt/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/seattlewaste2mqtt/master.svg)](http://codecov.io/github/mannkind/seattlewaste2mqtt?branch=master)

An experiment to publish Seattle Collection statuses/dates to MQTT.

## Use

The application can be locally built using `dotnet build` or you can utilize the multi-architecture Docker image(s).

### Example

```bash
docker run \
-e SEATTLEWASTE__SHARED__RESOURCES__0__Address="2133 N 61ST ST" \
-e SEATTLEWASTE__SHARED__RESOURCES__0__Slug="home" \
-e SEATTLEWASTE__SINK__BROKER="localhost" \
-e SEATTLEWASTE__SINK__DISCOVERYENABLED="true" \
mannkind/seattlewaste2mqtt:latest
```

OR

```bash
SEATTLEWASTE__SHARED__RESOURCES__0__Address="2133 N 61ST ST" \
SEATTLEWASTE__SHARED__RESOURCES__0__Slug="home" \
SEATTLEWASTE__SINK__BROKER="localhost" \
SEATTLEWASTE__SINK__DISCOVERYENABLED="true" \
./seattlewaste2mqtt 
```


## Configuration

Configuration happens via environmental variables

```bash
SEATTLEWASTE__SHARED__RESOURCES__#__Address      - The Address for a specific collection
SEATTLEWASTE__SHARED__RESOURCES__#__Slug         - The slug to identify the specific address
SEATTLEWASTE__SOURCE__POLLINGINTERVAL            - [OPTIONAL] The delay between collection lookups, defaults to "0.08:03:31"
SEATTLEWASTE__SINK__TOPICPREFIX                  - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/seattle_waste"
SEATTLEWASTE__SINK__DISCOVERYENABLED             - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
SEATTLEWASTE__SINK__DISCOVERYPREFIX              - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
SEATTLEWASTE__SINK__DISCOVERYNAME                - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "seattle_waste"
SEATTLEWASTE__SINK__BROKER                       - [OPTIONAL] The MQTT broker, defaults to "test.mosquitto.org"
SEATTLEWASTE__SINK__USERNAME                     - [OPTIONAL] The MQTT username, default to ""
SEATTLEWASTE__SINK__PASSWORD                     - [OPTIONAL] The MQTT password, default to ""
```
