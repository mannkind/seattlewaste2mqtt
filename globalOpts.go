package main

type globalOpts struct {
	Addresses sourceMapping `env:"SEATTLEWASTE_ADDRESS" envDefault:""`
}
