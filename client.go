package main

import (
	"fmt"
	"time"

	"github.com/mannkind/seattlewaste"
	log "github.com/sirupsen/logrus"
)

const (
	apiDateFormat       = "Mon, 2 Jan 2006"
	sensorTopicTemplate = "%s/%s/state"
	maxAPIAttempts      = 5
)

type client struct {
	observers map[observer]struct{}

	address        string
	alertWithin    time.Duration
	lookupInterval time.Duration
}

func newClient(config *config) *client {
	c := client{
		observers: map[observer]struct{}{},

		address:        config.Address,
		alertWithin:    config.AlertWithin,
		lookupInterval: config.LookupInterval,
	}

	return &c
}

func (c *client) run() {
	go c.loop(false)
}

func (c *client) register(l observer) {
	c.observers[l] = struct{}{}
}

func (c *client) publish(e event) {
	for o := range c.observers {
		o.receiveState(e)
	}
}

func (c *client) loop(once bool) {
	for {
		log.Info("Beginning lookup")
		now := time.Now()
		if info, err := c.lookup(now); err == nil {
			c.publish(event{
				version: 1,
				data:    c.adapt(info),
			})
		}
		log.Info("Ending lookup")

		if once {
			break
		}

		time.Sleep(c.lookupInterval)
	}
}

func (c *client) lookup(now time.Time) (seattlewaste.Collection, error) {
	none := seattlewaste.Collection{}
	swclient := seattlewaste.NewClient(c.address)

	localLoc, _ := time.LoadLocation("Local")
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 1, 0, localLoc)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, localLoc)
	lastTimestamp := firstOfMonth.Unix()
	todayTimestamp := today.Unix()
	apiCallCount := 0

	for lastTimestamp < todayTimestamp && apiCallCount <= maxAPIAttempts {
		results, err := swclient.GetCollections(lastTimestamp)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to fetch collection dates")
			return none, err
		}

		apiCallCount++

		if len(results) == 0 {
			return none, fmt.Errorf("No collection dates returned")
		}

		// Results from the 'web-service' do not always return as expected
		for _, result := range results {
			pTime, err := time.ParseInLocation(apiDateFormat, result.Start, localLoc)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Error parsing the datetime from the 'API'")
				continue
			}

			lastTimestamp = pTime.Unix()
			if lastTimestamp >= todayTimestamp {
				until := pTime.Sub(time.Now())
				result.Status = 0 <= until && until <= c.alertWithin
				return result, nil
			}
		}
	}

	return none, nil
}

func (c *client) adapt(info seattlewaste.Collection) eventData {
	return eventData{
		Start:            info.Start,
		Garbage:          info.Garbage,
		Recycling:        info.Recycling,
		FoodAndYardWaste: info.FoodAndYardWaste,
		Status:           info.Status,
	}
}
