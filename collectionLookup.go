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
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

const (
	apiDateFormat        = "Mon, 2 Jan 2006"
	apiAddressURL        = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCCAddress?pAddress=%s"
	apiCollectionDaysURL = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays?pAddress=%s&pApp=CC&start=%s"
	sensorTopicTemplate  = "%s/%s/state"
)

var binarySensors = []string{"Garbage", "Recycling", "FoodAndYardWaste", "Status"}
var stringSensors = []string{"Start"}

type apiResponse struct {
	Start            string
	Garbage          bool
	Recycling        bool
	FoodAndYardWaste bool
	Date             time.Time
	Status           bool
}

type collectionLookup struct {
	ClientID        string        `env:"MQTT_CLIENTID" envDefault:"DefaultSeattleWaste2MQTTClientID"`
	Broker          string        `env:"MQTT_BROKER" envDefault:"tcp://mosquitto.org:1883"`
	PubTopic        string        `env:"MQTT_PUBTOPIC" envDefault:"home/seattle_waste"`
	Discovery       bool          `env:"MQTT_DISCOVERY" envDefault:"false"`
	DiscoveryPrefix string        `env:"MQTT_DISCOVERYPREFIX" envDefault:"homeassistant"`
	DiscoveryName   string        `env:"MQTT_DISCOVERYNAME" envDefault:"seattle_waste"`
	Username        string        `env:"MQTT_USERNAME"`
	Password        string        `env:"MQTT_PASSWORD"`
	Address         string        `env:"SEATTLEWASTE_ADDRESS,required"`
	AlertWithin     time.Duration `env:"SEATTLEWASTE_ALERTWITHIN" envDefault:"24h"`
	LookupInterval  time.Duration `env:"SEATTLEWASTE_LOOKUPINTERVAL" envDefault:"8h"`

	client         mqtt.Client
	encodedAddress string
	lastPublished  string
}

// mqttDiscovery - Defines the fields used by this app for HA MQTT discovery
type mqttDiscovery struct {
	Name        string `json:"name"`
	StateTopic  string `json:"state_topic"`
	UniqueID    string `json:"unique_id,omitempty"`
	PayloadOn   string `json:"payload_on,omitempty"`
	PayloadOff  string `json:"payload_off,omitempty"`
	DeviceClass string `json:"device_class,omitempty"`
}

func (t *collectionLookup) start() error {
	log.Print("Connecting to MQTT")
	opts := mqtt.NewClientOptions().
		AddBroker(t.Broker).
		SetClientID(t.ClientID).
		SetOnConnectHandler(t.onConnect).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			log.Printf("Disconnected from MQTT: %s.", err)
		}).
		SetUsername(t.Username).
		SetPassword(t.Password)

	t.client = mqtt.NewClient(opts)
	if token := t.client.Connect(); !token.Wait() || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (t *collectionLookup) stop() {
	if t.client != nil && t.client.IsConnected() {
		t.client.Disconnect(0)
	}
}

func (t *collectionLookup) onConnect(client mqtt.Client) {
	log.Print("Connected to MQTT")

	if !client.IsConnected() {
		log.Print("Subscribe Error: Not Connected (Reloading Config?)")
		return
	}

	go t.loop(false)
}

func (t *collectionLookup) loop(once bool) {
	if t.Discovery {
		t.discovery()
	}

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

		time.Sleep(t.LookupInterval)
	}
}

func (t *collectionLookup) discovery() {
	sensorMap := map[string][]string{
		"binary_sensor": binarySensors,
		"sensor":        stringSensors,
	}
	for sensorType, sensors := range sensorMap {
		for _, sensor := range sensors {
			sensorSlug := strings.ToLower(sensor)
			mqd := mqttDiscovery{
				Name:       fmt.Sprintf("%s %s", t.DiscoveryName, sensor),
				StateTopic: fmt.Sprintf(sensorTopicTemplate, t.PubTopic, sensorSlug),
				UniqueID:   fmt.Sprintf("%s.%s", t.DiscoveryName, sensorSlug),
			}

			topic := fmt.Sprintf("%s/%s/%s/%s/config", t.DiscoveryPrefix, sensorType, t.DiscoveryName, sensorSlug)
			payloadBytes, _ := json.Marshal(mqd)
			payload := string(payloadBytes)

			t.publish(topic, payload)
		}
	}
}

func (t *collectionLookup) encodeAddress() error {
	// Only encode the address once
	if t.encodedAddress != "" {
		return nil
	}

	// GET the encoded adddress
	var body io.ReadCloser
	url := fmt.Sprintf(apiAddressURL, url.QueryEscape(t.Address))
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

func (t *collectionLookup) collectionLookup(now time.Time) (apiResponse, error) {
	noResult := apiResponse{}

	// Guard-clause for a blank encoded address
	if t.encodedAddress == "" {
		return noResult, errors.New("No encoded address found for collection lookup")
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 1, 0, time.UTC)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
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
			pTime, err := time.Parse(apiDateFormat, result.Start)
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

func (t *collectionLookup) publishCollectionInfo(info apiResponse) {
	until := info.Date.Sub(time.Now())
	info.Status = until >= 0 && until <= t.AlertWithin

	sensorMap := map[string][]string{
		"binary_sensor": binarySensors,
		"sensor":        stringSensors,
	}
	for sensorType, sensors := range sensorMap {
		for _, sensor := range sensors {
			sensorSlug := strings.ToLower(sensor)
			sensorValue := reflect.Indirect(reflect.ValueOf(info)).FieldByName(sensor)
			topic := fmt.Sprintf(sensorTopicTemplate, t.PubTopic, sensorSlug)
			payload := ""
			if sensorType == "binary_sensor" {
				payload = "OFF"
				if sensorValue.Bool() {
					payload = "ON"
				}
			} else if sensorType == "sensor" {
				payload = sensorValue.String()
			}

			t.publish(topic, payload)
		}
	}
}

func (t *collectionLookup) publish(topic string, payload string) {
	retain := true
	if token := t.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}
	t.lastPublished = fmt.Sprintf("Publishing - Topic:%s ; Payload: %s", topic, payload)
	log.Print(t.lastPublished)
}
