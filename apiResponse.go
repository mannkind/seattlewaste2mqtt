package main

import (
	"time"
)

type apiResponse struct {
	Start            string
	Garbage          bool
	Recycling        bool
	FoodAndYardWaste bool
	Date             time.Time
	Status           bool
}
