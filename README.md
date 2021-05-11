# seattlewaste2mqtt

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/seattlewaste2mqtt/blob/main/LICENSE.md)
[![Build Status](https://github.com/mannkind/seattlewaste2mqtt/workflows/Main%20Workflow/badge.svg)](https://github.com/mannkind/seattlewaste2mqtt/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/seattlewaste2mqtt/main.svg)](http://codecov.io/github/mannkind/seattlewaste2mqtt?branch=main)

An experiment to publish Seattle Collection statuses/dates to MQTT.

NOTE: As of 2021, I no longer live in the city of Seattle. I currently plan to keep this project active (unless it becomes a burden), but I might not notice issues immediately because I no longer rely on it.

## Use

The application can be locally built using `dotnet build` or you can utilize the multi-architecture Docker image(s).

### Example

```bash
docker run \
-e SEATTLEWASTE__RESOURCES__0__Address="2133 N 61ST ST" \
-e SEATTLEWASTE__RESOURCES__0__Slug="home" \
-e SEATTLEWASTE__MQTT__BROKER="localhost" \
-e SEATTLEWASTE__MQTT__DISCOVERYENABLED="true" \
mannkind/seattlewaste2mqtt:latest
```

OR

```bash
SEATTLEWASTE__RESOURCES__0__Address="2133 N 61ST ST" \
SEATTLEWASTE__RESOURCES__0__Slug="home" \
SEATTLEWASTE__MQTT__BROKER="localhost" \
SEATTLEWASTE__MQTT__DISCOVERYENABLED="true" \
./seattlewaste2mqtt 
```


## Configuration

Configuration happens via environmental variables

```bash
SEATTLEWASTE__RESOURCES__#__Address              - The Address for a specific collection
SEATTLEWASTE__RESOURCES__#__Slug                 - The slug to identify the specific address
SEATTLEWASTE__POLLINGINTERVAL                    - [OPTIONAL] The delay between collection lookups, defaults to "0.08:03:31"
SEATTLEWASTE__MQTT__TOPICPREFIX                  - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/seattle_waste"
SEATTLEWASTE__MQTT__DISCOVERYENABLED             - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
SEATTLEWASTE__MQTT__DISCOVERYPREFIX              - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
SEATTLEWASTE__MQTT__DISCOVERYNAME                - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "seattle_waste"
SEATTLEWASTE__MQTT__BROKER                       - [OPTIONAL] The MQTT broker, defaults to "test.mosquitto.org"
SEATTLEWASTE__MQTT__USERNAME                     - [OPTIONAL] The MQTT username, default to ""
SEATTLEWASTE__MQTT__PASSWORD                     - [OPTIONAL] The MQTT password, default to ""
```

## Prior Implementations

### Golang
* Last Commit: [efddb0703cb0e309c98f6f801ebf46da7ae12193](https://github.com/mannkind/seattlewaste2mqtt/commit/efddb0703cb0e309c98f6f801ebf46da7ae12193)
* Last Docker Image: [mannkind/seattlewaste2mqtt:v0.15.20055.0754](https://hub.docker.com/layers/mannkind/seattlewaste2mqtt/v0.15.20055.0754/images/sha256-6ad7368c88c46326e2ef755053885c113e35981081de38077ff73cf4d4ec08d4?context=explore)