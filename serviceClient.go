package main

import (
	"fmt"
	"time"

	"github.com/mannkind/seattlewaste"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const (
	apiDateFormat  = "Mon, 2 Jan 2006"
	maxAPIAttempts = 5
)

type serviceClient struct {
	serviceClientConfig
	stateUpdateChan stateChannel
}

func newServiceClient(serviceClientCfg serviceClientConfig, stateUpdateChan stateChannel) *serviceClient {
	c := serviceClient{
		serviceClientConfig: serviceClientCfg,
		stateUpdateChan:     stateUpdateChan,
	}

	log.WithFields(log.Fields{
		"SeattleWaste.Addresses":      c.Addresses,
		"SeattleWaste.AlertWithin":    c.AlertWithin,
		"SeattleWaste.LookupInterval": c.LookupInterval,
	}).Info("Service Client Environmental Settings")

	return &c
}

func (c *serviceClient) run() {
	// Run immediately
	go c.loop()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.LookupInterval), c.loop)
	sched.Start()
}

func (c *serviceClient) loop() {
	log.Info("Looping")
	now := time.Now()
	for address := range c.Addresses {
		info, err := c.lookup(address, now)
		if err != nil {
			continue
		}

		obj, err := c.adapt(address, info)
		if err != nil {
			continue
		}

		c.stateUpdateChan <- obj
	}

	log.WithFields(log.Fields{
		"sleep": c.LookupInterval,
	}).Info("Finished looping; sleeping")
}

func (c *serviceClient) lookup(address string, now time.Time) (seattlewaste.Collection, error) {
	log.WithFields(log.Fields{
		"address": address,
	}).Info("Looking up collection information for address")

	none := seattlewaste.Collection{}
	swclient := seattlewaste.NewClient(address)

	localLoc, _ := time.LoadLocation("Local")
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 1, 0, localLoc)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, localLoc)
	lastTimestamp := firstOfMonth.Unix()
	todayTimestamp := today.Unix()
	apiCallCount := 0

	for lastTimestamp < todayTimestamp && apiCallCount <= maxAPIAttempts {
		log.WithFields(log.Fields{
			"lastTimeStamp": lastTimestamp,
			"count":         apiCallCount,
		}).Debug("Sending API request(s)")

		results, err := swclient.GetCollections(lastTimestamp)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to fetch collection dates")
			return none, err
		}

		log.Debug("Finished sending API request(s)")

		apiCallCount++
		if len(results) == 0 {
			return none, fmt.Errorf("No collection dates returned")
		}

		// Results from the 'web-service' do not always return as expected
		for _, result := range results {
			log.WithFields(log.Fields{
				"result": result,
			}).Debug("Processing collection result")

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
				result.Status = 0 <= until && until <= c.AlertWithin

				log.Debug("Finished API request(s)")

				return result, nil
			}
		}
	}

	log.Debug("Finished API request(s); nothing found")
	return none, nil
}

func (c *serviceClient) adapt(address string, info seattlewaste.Collection) (collection, error) {
	log.WithFields(log.Fields{
		"address":    address,
		"collection": info,
	}).Debug("Adapting collection information")

	obj := collection{
		Address:          address,
		Start:            info.Start,
		Garbage:          info.Garbage,
		Recycling:        info.Recycling,
		FoodAndYardWaste: info.FoodAndYardWaste,
		Status:           info.Status,
	}

	log.Debug("Finished adapting collection information")
	return obj, nil
}
