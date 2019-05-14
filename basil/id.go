package basil

import (
	"regexp"
)

// IDRegExpPattern is the regular expression for a valid identifier
const IDRegExpPattern = "[a-z][a-z0-9]*(?:_[a-z0-9]+)*"

// IDRegExp is a compiled regular expression object for a valid identifier
var IDRegExp = regexp.MustCompile("^" + IDRegExpPattern + "$")

// Keywords are reserved strings and may not be used as identifiers.
var Keywords = []string{"true", "false", "nil", "map"}

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
}
