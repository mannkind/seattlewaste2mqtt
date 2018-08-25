package main

import (
	"log"

	"github.com/mannkind/seattle_waste_mqtt/cmd"
)

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

func main() {
	log.Printf("Seattle Waste Version: %s", Version)
	cmd.Execute()
}
