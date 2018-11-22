//+build wireinject

package main

import (
	"github.com/google/go-cloud/wire"
)

// InitializeCollectionLookup - Compile-time DI
func InitializeCollectionLookup() *CollectionLookup {
	wire.Build(NewConfig, NewMQTTFuncWrapper, NewCollectionLookup)

	return &CollectionLookup{}
}
