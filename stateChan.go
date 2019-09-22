package main

type stateChannel = chan collection

func newStateChannel() stateChannel {
	return make(stateChannel, 100)
}
