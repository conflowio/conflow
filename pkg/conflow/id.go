// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

// Keywords are reserved strings and may not be used as identifiers.
var Keywords = []string{}

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
//
//counterfeiter:generate . Identifiable
type Identifiable interface {
	ID() ID
}
