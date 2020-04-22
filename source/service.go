package source

import (
	"fmt"
	"time"

	"github.com/mannkind/seattlewaste2mqtt/lib"
	log "github.com/sirupsen/logrus"
)

// Service is for reading a directly from a source system
type Service struct {
}

// NewService creates a new Service for reading a directly from a source system
func NewService() *Service {
	c := Service{}

	return &c
}

func (c *Service) lookup(address string, now time.Time, alertWithin time.Duration) (*lib.Collection, error) {
	log.WithFields(log.Fields{
		"address": address,
	}).Info("Looking up collection information for address")

	swclient := lib.NewClient(address)

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
			return nil, err
		}

		log.Debug("Finished sending API request(s)")

		apiCallCount++
		if len(results) == 0 {
			return nil, fmt.Errorf("No collection dates returned")
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
				result.Status = 0 <= until && until <= alertWithin

				log.Debug("Finished API request(s)")

				return &result, nil
			}
		}
	}

	log.WithFields(log.Fields{
		"address": address,
	}).Info("Finished API request(s); nothing found")
	return nil, fmt.Errorf("Terrible, horrible, no good, very bad things")
}
