package ocl

// Identifiable makes an object to have a string identifier
type Identifiable interface {
	ID() string
}

// IDRegistry provides information about existing identifiers and able to generate new ones
type IDRegistry interface {
	IDExists(string) bool
	GenerateID() string
	RegisterID(string) error
}

// IDRegistryAware defines an interface to retrieve an identifer registry
type IDRegistryAware interface {
	GetIDRegistry() IDRegistry
}
