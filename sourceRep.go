package main

type sourceRep struct {
	Address          string `mqttDiscoveryType:",ignore"`
	Start            string `mqttDiscoveryType:"sensor"`
	Garbage          bool   `mqttDiscoveryType:"binary_sensor"`
	Recycling        bool   `mqttDiscoveryType:"binary_sensor"`
	FoodAndYardWaste bool   `mqttDiscoveryType:"binary_sensor"`
	Status           bool   `mqttDiscoveryType:"binary_sensor"`
}
