package main

type eventData struct {
	Start            string
	Garbage          bool
	Recycling        bool
	FoodAndYardWaste bool
	Status           bool
}

type event struct {
	version int64
	data    eventData
}

type observer interface {
	receive(event)
}

type publisher interface {
	register(observer)
	publish(event)
}
