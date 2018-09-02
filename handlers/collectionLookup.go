package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

const (
	apiDateFormat        = "Mon, 2 Jan 2006"
	apiAddressURL        = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCCAddress?pAddress=%s"
	apiCollectionDaysURL = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays?pAddress=%s&pApp=CC&start=%s"
)

type apiResponse struct {
	Start            string
	Garbage          bool
	Recycling        bool
	FoodAndYardWaste bool
	Date             time.Time
	Status           string
}

// CollectionLookup - MQTT all the things!
type CollectionLookup struct {
	Client   mqtt.Client
	Settings struct {
		ClientID string
		Broker   string
		SubTopic string
		PubTopic string
		Username string
		Password string
	}
	Control struct {
		Address        string
		EncodedAddress string
		AlertWithin    time.Duration
		LookupInterval time.Duration
	}
	LastPublished string
}

// Start - Connect and Subscribe
func (t *CollectionLookup) Start() error {
	log.Println("Connecting to MQTT: ", t.Settings.Broker)
	opts := mqtt.NewClientOptions().
		AddBroker(t.Settings.Broker).
		SetClientID(t.Settings.ClientID).
		SetOnConnectHandler(t.onConnect).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			log.Printf("Disconnected from MQTT: %s.", err)
		}).
		SetUsername(t.Settings.Username).
		SetPassword(t.Settings.Password)

	t.Client = mqtt.NewClient(opts)
	if token := t.Client.Connect(); !token.Wait() || token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Stop - Disconnect
func (t *CollectionLookup) Stop() {
	if t.Client != nil && t.Client.IsConnected() {
		t.Client.Disconnect(0)
	}
}

func (t *CollectionLookup) onConnect(client mqtt.Client) {
	log.Println("Connected to MQTT")

	if !client.IsConnected() {
		log.Print("Subscribe Error: Not Connected (Reloading Config?)")
		return
	}

	go t.loop()
}

func (t *CollectionLookup) loop() {
	for {
		log.Println("Beginning address encoding")
		t.encodeAddress()
		log.Println("Ending address encoding")

		log.Println("Beginning collection lookup")
		if collectionInfo, err := t.collectionLookup(); collectionInfo.Start != "" && err == nil {
			t.publishCollectionInfo(collectionInfo)
		} else {
			log.Println(err)
		}
		log.Println("Ending collection lookup")

		time.Sleep(t.Control.LookupInterval)
	}
}

func (t *CollectionLookup) encodeAddress() error {
	// Only encode the address once
	if t.Control.EncodedAddress != "" {
		return nil
	}

	// GET the encoded adddress
	var body io.ReadCloser
	url := fmt.Sprintf(apiAddressURL, url.QueryEscape(t.Control.Address))
	if resp, err := http.Get(url); err == nil && resp.StatusCode == http.StatusOK {
		body = resp.Body
	} else {
		log.Println(err)
		return errors.New("Unble to encode the address")
	}

	// Decode the response
	var result []string
	json.NewDecoder(body).Decode(&result)

	// Store the result
	if len(result) > 0 {
		t.Control.EncodedAddress = result[0]
	}

	return nil
}

func (t *CollectionLookup) collectionLookup() (apiResponse, error) {
	noResult := apiResponse{}

	// Guard-clause for a blank encoded address
	if t.Control.EncodedAddress == "" {
		return noResult, errors.New("No encoded address found for collection lookup")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 1, 0, time.UTC)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastTimestamp := firstOfMonth.Unix()
	todayTimestamp := today.Unix()

	var collectionInfo apiResponse
	for lastTimestamp < todayTimestamp {
		encodedAddress := url.QueryEscape(t.Control.EncodedAddress)
		timeCheck := url.QueryEscape(fmt.Sprintf("%d", lastTimestamp))

		// Get the collection days
		var body io.ReadCloser
		url := fmt.Sprintf(apiCollectionDaysURL, encodedAddress, timeCheck)
		if resp, err := http.Get(url); err == nil && resp.StatusCode == http.StatusOK {
			body = resp.Body
		} else {
			log.Println(err)
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
				log.Println(err)
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
	info.Status = "OFF"
	if until >= 0 && until <= t.Control.AlertWithin {
		info.Status = "ON"
	}

	// Publish the attributes about the waste to be picked up
	attrBytes, err := json.Marshal(info)
	if err != nil {
		return
	}

	t.publish(t.Client, t.Settings.PubTopic, string(attrBytes))
}

func (t *CollectionLookup) publish(client mqtt.Client, topic string, payload string) {
	retain := true
	if token := client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}
	t.LastPublished = fmt.Sprintf("Publishing - Topic:%s ; Payload: %s", topic, payload)
	log.Println(t.LastPublished)
}
