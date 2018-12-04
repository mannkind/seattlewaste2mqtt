//+build wireinject

package main

import (
	"github.com/google/wire"
)

// InitializeCollectionLookup - Compile-time DI
func InitializeCollectionLookup() *CollectionLookup {
	wire.Build(NewConfig, NewMQTTFuncWrapper, NewCollectionLookup)

	return &CollectionLookup{}
}
