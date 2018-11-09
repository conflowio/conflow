package basil

import "github.com/opsidian/basil/variable"

// Identifiable makes an object to have a string identifier
type Identifiable interface {
	ID() variable.ID
}

// IDRegistry provides information about existing identifiers and able to generate new ones
type IDRegistry interface {
	IDExists(variable.ID) bool
	GenerateID() variable.ID
	RegisterID(variable.ID) error
}

// IDRegistryAware defines an interface to retrieve an identifer registry
type IDRegistryAware interface {
	IDRegistry() IDRegistry
}
