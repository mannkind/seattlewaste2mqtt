package main

type globalClientConfig struct {
	Addresses addressMapping `env:"SEATTLEWASTE_ADDRESS" envDefault:""`
}
