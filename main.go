package main

import (
	"log"
)

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Printf("Seattle Waste Version: %s", Version)

	cl := InitializeCollectionLookup()
	if err := cl.Run(); err != nil {
		log.Panicf("Error starting collection lookup process: %s", err)
	}

	select {}
}
