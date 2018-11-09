package main

import (
	"log"

	"github.com/caarlos0/env"
)

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Printf("Seattle Waste Version: %s", Version)
	log.Print("Starting Process")

	controller := collectionLookup{}
	if err := env.Parse(&controller); err != nil {
		log.Panicf("Error unmarshaling configuration: %s", err)
	}

	redactedPassword := ""
	if len(controller.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.Printf("Environmental Settings:")
	log.Printf("  * ClientID      : %s", controller.ClientID)
	log.Printf("  * Broker        : %s", controller.Broker)
	log.Printf("  * PubTopic      : %s", controller.PubTopic)
	log.Printf("  * Username      : %s", controller.Username)
	log.Printf("  * Password      : %s", redactedPassword)
	log.Printf("  * Address       : %s", controller.Address)
	log.Printf("  * AlertWithin   : %s", controller.AlertWithin)
	log.Printf("  * LookupInterval: %s", controller.LookupInterval)

	if err := controller.start(); err != nil {
		log.Panicf("Error starting collection lookup handler: %s", err)
	}

	// log.Print("Ending Process")
	select {}
}
