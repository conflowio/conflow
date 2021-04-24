// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"regexp"
)

// IDRegExpPattern is the regular expression for a valid identifier
const IDRegExpPattern = "[a-z][a-z0-9]*(?:_[a-z0-9]+)*"

// IDRegExp is a compiled regular expression object for a valid identifier
var IDRegExp = regexp.MustCompile("^" + IDRegExpPattern + "$")

// Keywords are reserved strings and may not be used as identifiers.
var Keywords = []string{
	"map",
}

// ID contains an identifier
type ID string

// String returns with the ID string
func (i ID) String() string {
	return string(i)
}

func (i ID) Concat(s string) ID {
	return ID(string(i) + s)
}

// Identifiable makes an object to have a string identifier and have an identifiable parent
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Identifiable
type Identifiable interface {
	ID() ID
}
