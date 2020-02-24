package source

import (
	"fmt"
	"time"

	"github.com/mannkind/seattlewaste"
	"github.com/mannkind/seattlewaste2mqtt/shared"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const (
	apiDateFormat  = "Mon, 2 Jan 2006"
	maxAPIAttempts = 5
)

// Reader is for reading a shared representation out of a source system
type Reader struct {
	opts     Opts
	outgoing chan<- shared.Representation
	service  *Service
}

// NewReader creates a new Reader for reading a shared representation out of a source system
func NewReader(opts Opts, outgoing chan<- shared.Representation, service *Service) *Reader {
	c := Reader{
		opts:     opts,
		outgoing: outgoing,
		service:  service,
	}

	return &c
}

// Run starts the Reader
func (c *Reader) Run() {
	// Log service settings
	c.logSettings()

	// Run immediately
	c.poll()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.opts.LookupInterval), c.poll)
	sched.Start()
}

func (c *Reader) logSettings() {
	log.WithFields(log.Fields{
		"SeattleWaste.Addresses":      c.opts.Addresses,
		"SeattleWaste.AlertWithin":    c.opts.AlertWithin,
		"SeattleWaste.LookupInterval": c.opts.LookupInterval,
	}).Info("Service Client Environmental Settings")
}

func (c *Reader) poll() {
	log.Info("Polling")
	now := time.Now()
	for address := range c.opts.Addresses {
		info, err := c.service.lookup(address, now, c.opts.AlertWithin)
		if err != nil {
			continue
		}

		c.outgoing <- c.adapt(address, info)
	}

	log.WithFields(log.Fields{
		"sleep": c.opts.LookupInterval,
	}).Info("Finished polling; sleeping")
}

func (c *Reader) adapt(address string, info *seattlewaste.Collection) shared.Representation {
	obj := shared.Representation{
		Address:          address,
		Start:            info.Start,
		Garbage:          info.Garbage,
		Recycling:        info.Recycling,
		FoodAndYardWaste: info.FoodAndYardWaste,
		Status:           info.Status,
	}

	return obj
}
