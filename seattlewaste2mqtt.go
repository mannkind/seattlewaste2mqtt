package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
	mqttExtHA "github.com/mannkind/paho.mqtt.golang.ext/ha"
)

const (
	apiDateFormat        = "Mon, 2 Jan 2006"
	apiAddressURL        = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCCAddress?pAddress=%s"
	apiCollectionDaysURL = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays?pAddress=%s&pApp=CC&Start=%s"
	sensorTopicTemplate  = "%s/%s/state"
)

var binarySensors = []string{"Garbage", "Recycling", "FoodAndYardWaste", "Status"}
var stringSensors = []string{"Start"}

// SeattleWaste2Mqtt - Lookup collection information on seattle.gov.
type SeattleWaste2Mqtt struct {
	discovery       bool
	discoveryPrefix string
	discoveryName   string
	topicPrefix     string
	address         string
	alertWithin     time.Duration
	lookupInterval  time.Duration

	client         mqtt.Client
	encodedAddress string
}

// NewSeattleWaste2Mqtt - Returns a new reference to a fully configured object.
func NewSeattleWaste2Mqtt(config *Config, mqttFuncWrapper *mqttExtDI.MQTTFuncWrapper) *SeattleWaste2Mqtt {
	cl := SeattleWaste2Mqtt{
		discovery:       config.MQTT.Discovery,
		discoveryPrefix: config.MQTT.DiscoveryPrefix,
		discoveryName:   config.MQTT.DiscoveryName,
		topicPrefix:     config.MQTT.TopicPrefix,
		address:         config.Address,
		alertWithin:     config.AlertWithin,
		lookupInterval:  config.LookupInterval,
	}

	opts := mqttFuncWrapper.
		ClientOptsFunc().
		AddBroker(config.MQTT.Broker).
		SetClientID(config.MQTT.ClientID).
		SetOnConnectHandler(cl.onConnect).
		SetConnectionLostHandler(cl.onDisconnect).
		SetUsername(config.MQTT.Username).
		SetPassword(config.MQTT.Password)

	cl.client = mqttFuncWrapper.ClientFunc(opts)

	return &cl
}

// Run - Start the collection lookup process
func (t *SeattleWaste2Mqtt) Run() error {
	log.Print("Connecting to MQTT")
	if token := t.client.Connect(); !token.Wait() || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (t *SeattleWaste2Mqtt) onConnect(client mqtt.Client) {
	log.Print("Connected to MQTT")

	if !client.IsConnected() {
		log.Print("Subscribe Error: Not Connected (Reloading Config?)")
		return
	}

	if t.discovery {
		t.publishDiscovery()
	}

	go t.loop(false)
}

func (t *SeattleWaste2Mqtt) onDisconnect(client mqtt.Client, err error) {
	log.Printf("Disconnected from MQTT: %s.", err)
}

func (t *SeattleWaste2Mqtt) publishDiscovery() {
	obj := reflect.ValueOf(apiResponse{})
	for i := 0; i < obj.NumField(); i++ {
		sensor := strings.ToLower(obj.Type().Field(i).Name)
		val := obj.Field(i)
		sensorType := ""

		switch val.Kind() {
		case reflect.Bool:
			sensorType = "binary_sensor"
		case reflect.String:
			sensorType = "sensor"
		}

		if sensorType == "" {
			continue
		}

		mqd := mqttExtHA.MQTTDiscovery{
			DiscoveryPrefix: t.discoveryPrefix,
			Component:       sensorType,
			NodeID:          t.discoveryName,
			ObjectID:        sensor,
			Name:            fmt.Sprintf("%s %s", t.discoveryName, sensor),
			StateTopic:      fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, sensor),
			UniqueID:        fmt.Sprintf("%s.%s", t.discoveryName, sensor),
		}

		mqd.PublishDiscovery(t.client)
	}
}

func (t *SeattleWaste2Mqtt) loop(once bool) {
	for {
		log.Print("Beginning address encoding")
		t.encodeAddress()
		log.Print("Ending address encoding")

		log.Print("Beginning collection lookup")
		now := time.Now()
		if collectionInfo, err := t.collectionLookup(now); collectionInfo.Start != "" && err == nil {
			t.publishCollectionInfo(collectionInfo)
		} else {
			log.Print(err)
		}
		log.Print("Ending collection lookup")

		if once {
			break
		}

		time.Sleep(t.lookupInterval)
	}
}

func (t *SeattleWaste2Mqtt) encodeAddress() error {
	// Only encode the address once
	if t.encodedAddress != "" {
		return nil
	}

	// GET the encoded adddress
	url := fmt.Sprintf(apiAddressURL, url.QueryEscape(t.address))
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Print(err)
		return errors.New("Unble to encode the address")
	}

	// Decode the response
	var result []string
	json.NewDecoder(resp.Body).Decode(&result)

	// Store the result
	if len(result) > 0 {
		t.encodedAddress = result[0]
	}

	return nil
}

func (t *SeattleWaste2Mqtt) collectionLookup(now time.Time) (apiResponse, error) {
	noResult := apiResponse{}

	// Guard-clause for a blank encoded address
	if t.encodedAddress == "" {
		return noResult, errors.New("No encoded address found for collection lookup")
	}

	localLoc, _ := time.LoadLocation("Local")
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 1, 0, localLoc)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, localLoc)
	lastTimestamp := firstOfMonth.Unix()
	todayTimestamp := today.Unix()

	var collectionInfo apiResponse
	for lastTimestamp < todayTimestamp {
		encodedAddress := url.QueryEscape(t.encodedAddress)
		timeCheck := url.QueryEscape(fmt.Sprintf("%d", lastTimestamp))

		// Get the collection days
		url := fmt.Sprintf(apiCollectionDaysURL, encodedAddress, timeCheck)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Print(err)
			return noResult, errors.New("Unable to fetch collection dates")
		}

		var results []apiResponse
		json.NewDecoder(resp.Body).Decode(&results)

		if len(results) == 0 {
			return noResult, errors.New("No collection dates returned")
		}

		// Results from the 'web-service' do not always return as expected
		for _, result := range results {
			pTime, err := time.ParseInLocation(apiDateFormat, result.Start, localLoc)
			if err != nil {
				log.Print(err)
				continue
			}

			lastTimestamp = pTime.Unix()
			if lastTimestamp >= todayTimestamp {
				collectionInfo = result
				collectionInfo.Date = pTime
				break
			}
		}
	}

	return collectionInfo, nil
}

func (t *SeattleWaste2Mqtt) publishCollectionInfo(info apiResponse) {
	until := info.Date.Sub(time.Now())
	info.Status = until >= 0 && until <= t.alertWithin

	obj := reflect.ValueOf(info)
	for i := 0; i < obj.NumField(); i++ {
		sensor := strings.ToLower(obj.Type().Field(i).Name)
		val := obj.Field(i)

		topic := fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, sensor)
		payload := ""

		switch val.Kind() {
		case reflect.Bool:
			payload = "OFF"
			if val.Bool() {
				payload = "ON"
			}
		case reflect.String:
			payload = val.String()
		}

		if payload == "" {
			continue
		}

		t.publish(topic, payload)
	}
}

func (t *SeattleWaste2Mqtt) publish(topic string, payload string) {
	retain := true
	if token := t.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}

	log.Print(fmt.Sprintf("Publishing - Topic: %s ; Payload: %s", topic, payload))
}
