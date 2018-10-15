// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package variable

import "errors"

// ErrNotFound error is used when the given variable does not exist
var ErrNotFound = errors.New("Variable not found")

// Provider is an interface for looking up variables
//go:generate counterfeiter . Provider
type Provider interface {
	GetVar(name string) (interface{}, bool)
	LookupVar(lookup LookUp) (interface{}, error)
}

// LookUp is a variable lookup function
type LookUp func(provider Provider) (interface{}, error)
