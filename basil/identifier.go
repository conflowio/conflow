package basil

// ID contains a valid identifier
type ID string

// String returns with the ID string
func (i ID) String() string {
	return string(i)
}

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
