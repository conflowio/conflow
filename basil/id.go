package basil

import (
	"regexp"

	"github.com/opsidian/parsley/parsley"
)

// IDRegExpPattern is the regular expression for a valid identifier
const IDRegExpPattern = "[a-z][a-z0-9]*(?:_[a-z0-9]+)*"

// IDRegExp is a compiled regular expression object for a valid identifier
var IDRegExp = regexp.MustCompile("^" + IDRegExpPattern + "$")

// ID contains an identifier
type ID string

// String returns with the ID string
func (i ID) String() string {
	return string(i)
}

// Identifiable makes an object to have a string identifier and have an identifiable parent
//go:generate counterfeiter . Identifiable
type Identifiable interface {
	ID() ID
	ParentID() ID
}

// IdentifiableNode defines a node which contains an identifier
//go:generate counterfeiter . IdentifiableNode
type IdentifiableNode interface {
	parsley.Node
	Identifiable
}

// IDRegistry provides information about existing identifiers and able to generate new ones
type IDRegistry interface {
	IDExists(ID) bool
	GenerateID() ID
	RegisterID(ID) error
}

// IDRegistryAware defines an interface to retrieve an identifier registry
type IDRegistryAware interface {
	IDRegistry() IDRegistry
}
