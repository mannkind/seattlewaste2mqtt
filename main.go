package main

import (
	"log"

	"go.uber.org/dig"
)

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Printf("Seattle Waste Version: %s", Version)

	c := buildContainer()
	err := c.Invoke(func(cl *CollectionLookup) error {
		return cl.Run()
	})

	if err != nil {
		log.Panicf("Error starting collection lookup process: %s", err)
	}

	select {}
}

func buildContainer() *dig.Container {
	c := dig.New()
	c.Provide(NewConfig)
	c.Provide(NewMQTTFuncWrapper)
	c.Provide(NewCollectionLookup)

	return c
}
