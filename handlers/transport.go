package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

const (
	seattleWasteDateFormat = "Mon, 2 Jan 2006"
)

type seattleWasteJSONResponse struct {
	Start            string
	Garbage          bool
	Recycling        bool
	FoodAndYardWaste bool
	Date             time.Time
    Status           string
}

// SeattleWasteMQTT - MQTT all the things!
type SeattleWasteMQTT struct {
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
		Address         string
		EncodedAddress  string
		AlertWithin time.Duration
		LookupInterval time.Duration
	}
	LastPublished string
}

// Start - Connect and Subscribe
func (t *SeattleWasteMQTT) Start() error {
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
func (t *SeattleWasteMQTT) Stop() {
	if t.Client != nil && t.Client.IsConnected() {
		t.Client.Disconnect(0)
	}
}

func (t *SeattleWasteMQTT) onConnect(client mqtt.Client) {
	log.Println("Connected to MQTT")

	if !client.IsConnected() {
		log.Print("Subscribe Error: Not Connected (Reloading Config?)")
		return
	}

	lookup := func() {
		log.Println("Beginning collection lookup")
		t.encodeAddress()
		collectionInfo, err := t.collectionLookup()
		if err == nil {
			t.publishCollectionInfo(collectionInfo)
		} else {
            log.Println(err)
		}
		log.Println("Ending collection lookup")
	}

	for {
		lookup()
		time.Sleep(t.Control.LookupInterval)
	}
}

func (t *SeattleWasteMQTT) encodeAddress() {
	if t.Control.EncodedAddress != "" {
		return
	}

	var result []string
	url := fmt.Sprintf("https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCCAddress?pAddress=%s", url.QueryEscape(t.Control.Address))
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	json.NewDecoder(resp.Body).Decode(&result)
	if len(result) > 0 {
		t.Control.EncodedAddress = result[0]
	}
}

func (t *SeattleWasteMQTT) collectionLookup() (seattleWasteJSONResponse, error) {
	noResult := seattleWasteJSONResponse{}
	if t.Control.EncodedAddress == "" {
		return noResult, errors.New("No encoded address found for collection lookup")
	}

	var collectionInfo seattleWasteJSONResponse
	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	currentLocation := now.Location()
	today := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	lastCheck := firstOfMonth.Unix()
	todayCmp := today.Unix()

	for lastCheck < todayCmp {
		var results []seattleWasteJSONResponse

		encodedAddress := url.QueryEscape(t.Control.EncodedAddress)
		timeCheck := url.QueryEscape(fmt.Sprintf("%d", lastCheck))
		url := fmt.Sprintf("https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays?pAddress=%s&pApp=CC&start=%s", encodedAddress, timeCheck)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println(err)
			return noResult, errors.New("Unable to fetch collection dates")
		}

		json.NewDecoder(resp.Body).Decode(&results)

		if len(results) == 0 {
			return noResult, errors.New("No collection dates returned")
		}

		// Results from the 'web-service' do not always return as expected
		for _, result := range results {
			pTime, err := time.Parse(seattleWasteDateFormat, result.Start)
			if err != nil {
				log.Println(err)
				continue
			}

			lastCheck = pTime.Unix()
			if pTime.Unix() >= todayCmp {
				collectionInfo = result
				collectionInfo.Date = pTime
				break
			}
		}
	}

	return collectionInfo, nil
}

func (t *SeattleWasteMQTT) publishCollectionInfo(info seattleWasteJSONResponse) {
	info.Status = "OFF"
    alertDate := info.Date.Add(-1 * t.Control.AlertWithin)
    now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
    today := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, time.UTC)
    if alertDate.Equal(today) {
        info.Status = "ON"
    }

    // Publish the attributes about the waste to be picked up
	attrBytes, err := json.Marshal(info)
	if err != nil {
		return
	}

	t.publish(t.Client, t.Settings.PubTopic, string(attrBytes))
}

func (t *SeattleWasteMQTT) publish(client mqtt.Client, topic string, payload string) {
    retain := true
	if token := client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}
	t.LastPublished = fmt.Sprintf("%s %s", topic, payload)
}
