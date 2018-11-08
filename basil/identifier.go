package basil

// ID is the identifier type
type ID string

// Identifiable makes an object to have a string identifier
type Identifiable interface {
	ID() ID
}

// IDRegistry provides information about existing identifiers and able to generate new ones
type IDRegistry interface {
	IDExists(ID) bool
	GenerateID() ID
	RegisterID(ID) error
}

// IDRegistryAware defines an interface to retrieve an identifer registry
type IDRegistryAware interface {
	IDRegistry() IDRegistry
}
