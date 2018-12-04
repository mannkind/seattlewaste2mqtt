package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	apiDateFormat        = "Mon, 2 Jan 2006"
	apiAddressURL        = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCCAddress?pAddress=%s"
	apiCollectionDaysURL = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays?pAddress=%s&pApp=CC&Start=%s"
	sensorTopicTemplate  = "%s/%s/state"
)

var binarySensors = []string{"Garbage", "Recycling", "FoodAndYardWaste", "Status"}
var stringSensors = []string{"Start"}

// CollectionLookup - Lookup collection information on seattle.gov.
type CollectionLookup struct {
	discovery       bool
	discoveryPrefix string
	discoveryName   string
	topicPrefix     string
	address         string
	alertWithin     time.Duration
	lookupInterval  time.Duration

	client         mqtt.Client
	encodedAddress string
	lastPublished  string
}

// NewCollectionLookup - Returns a new, configured CollectionLoookup object.
func NewCollectionLookup(config *Config, mqttFuncWrapper *MQTTFuncWrapper) *CollectionLookup {
	cl := CollectionLookup{
		discovery:       config.Discovery,
		discoveryPrefix: config.DiscoveryPrefix,
		discoveryName:   config.DiscoveryName,
		topicPrefix:     config.TopicPrefix,
		address:         config.Address,
		alertWithin:     config.AlertWithin,
		lookupInterval:  config.LookupInterval,
	}

	opts := mqttFuncWrapper.
		clientOptsFunc().
		AddBroker(config.Broker).
		SetClientID(config.ClientID).
		SetOnConnectHandler(cl.onConnect).
		SetConnectionLostHandler(cl.onDisconnect).
		SetUsername(config.Username).
		SetPassword(config.Password)

	cl.client = mqttFuncWrapper.clientFunc(opts)

	return &cl
}

// Run - Start the collection lookup process
func (t *CollectionLookup) Run() error {
	log.Print("Connecting to MQTT")
	if token := t.client.Connect(); !token.Wait() || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (t *CollectionLookup) onConnect(client mqtt.Client) {
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

func (t *CollectionLookup) onDisconnect(client mqtt.Client, err error) {
	log.Printf("Disconnected from MQTT: %s.", err)
}

func (t *CollectionLookup) publishDiscovery() {
	sensorMap := map[string][]string{
		"binary_sensor": binarySensors,
		"sensor":        stringSensors,
	}
	for sensorType, sensors := range sensorMap {
		for _, sensor := range sensors {
			sensorSlug := strings.ToLower(sensor)
			mqd := mqttDiscovery{
				Name:       fmt.Sprintf("%s %s", t.discoveryName, sensor),
				StateTopic: fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, sensorSlug),
				UniqueID:   fmt.Sprintf("%s.%s", t.discoveryName, sensorSlug),
			}

			topic := fmt.Sprintf("%s/%s/%s/%s/config", t.discoveryPrefix, sensorType, t.discoveryName, sensorSlug)
			payloadBytes, _ := json.Marshal(mqd)
			payload := string(payloadBytes)

			t.publish(topic, payload)
		}
	}
}

func (t *CollectionLookup) loop(once bool) {
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

func (t *CollectionLookup) encodeAddress() error {
	// Only encode the address once
	if t.encodedAddress != "" {
		return nil
	}

	// GET the encoded adddress
	var body io.ReadCloser
	url := fmt.Sprintf(apiAddressURL, url.QueryEscape(t.address))
	if resp, err := http.Get(url); err == nil && resp.StatusCode == http.StatusOK {
		body = resp.Body
	} else {
		log.Print(err)
		return errors.New("Unble to encode the address")
	}

	// Decode the response
	var result []string
	json.NewDecoder(body).Decode(&result)

	// Store the result
	if len(result) > 0 {
		t.encodedAddress = result[0]
	}

	return nil
}

func (t *CollectionLookup) collectionLookup(now time.Time) (apiResponse, error) {
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
		var body io.ReadCloser
		url := fmt.Sprintf(apiCollectionDaysURL, encodedAddress, timeCheck)
		if resp, err := http.Get(url); err == nil && resp.StatusCode == http.StatusOK {
			body = resp.Body
		} else {
			log.Print(err)
			return noResult, errors.New("Unable to fetch collection dates")
		}

		var results []apiResponse
		json.NewDecoder(body).Decode(&results)

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

func (t *CollectionLookup) publishCollectionInfo(info apiResponse) {
	until := info.Date.Sub(time.Now())
	info.Status = until >= 0 && until <= t.alertWithin

	sensorMap := map[string][]string{
		"binary_sensor": binarySensors,
		"sensor":        stringSensors,
	}
	for _, sensors := range sensorMap {
		for _, sensor := range sensors {
			sensorSlug := strings.ToLower(sensor)
			sensorValue := reflect.Indirect(reflect.ValueOf(info)).FieldByName(sensor)
			topic := fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, sensorSlug)
			payload := ""
			switch sensorValue.Kind() {
			case reflect.Bool:
				payload = "OFF"
				if sensorValue.Bool() {
					payload = "ON"
				}
			case reflect.Int:
				payload = strconv.FormatInt(sensorValue.Int(), 10)
			case reflect.String:
				payload = sensorValue.String()
			}

			t.publish(topic, payload)
		}
	}
}

func (t *CollectionLookup) publish(topic string, payload string) {
	retain := true
	if token := t.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}
	t.lastPublished = fmt.Sprintf("Publishing - Topic:%s ; Payload: %s", topic, payload)
	log.Print(t.lastPublished)
}
