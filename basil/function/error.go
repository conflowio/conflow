// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import "fmt"

// Error is a function error which can contain an argument number
type Error struct {
	Err      error
	ArgIndex int
}

// NewError creates a new function error
func NewError(argIndex int, err error) *Error {
	return &Error{
		ArgIndex: argIndex,
		Err:      err,
	}
}

// NewErrorf creates a new function error
func NewErrorf(argIndex int, format string, args ...interface{}) *Error {
	return &Error{
		ArgIndex: argIndex,
		Err:      fmt.Errorf(format, args...),
	}
}

// Error returns with the error message
func (e *Error) Error() string {
	return e.Err.Error()
}
