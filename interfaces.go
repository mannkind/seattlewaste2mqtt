package main

import "github.com/mannkind/seattlewaste"

type event struct {
	version int64
	data    seattlewaste.Collection
}

type observer interface {
	receive(event)
}

type publisher interface {
	register(observer)
	publish(event)
}
