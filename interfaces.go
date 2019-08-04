package main

type eventData struct {
	Start            string `mqttDiscoveryType:"sensor"`
	Garbage          bool   `mqttDiscoveryType:"binary_sensor"`
	Recycling        bool   `mqttDiscoveryType:"binary_sensor"`
	FoodAndYardWaste bool   `mqttDiscoveryType:"binary_sensor"`
	Status           bool   `mqttDiscoveryType:"binary_sensor"`
}

type event struct {
	version int64
	data    eventData
}

type observer interface {
	receiveState(event)
	receiveCommand(int64, event)
}

type publisher interface {
	register(observer)
}
