// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package multierror

func Append(e1 error, e2 error) error {
	switch {
	case e1 == nil:
		return e2
	case e2 == nil:
		return e1
	}

	if e, ok := e1.(*Error); ok {
		return e.Append(e2)
	}

	e := &Error{errors: []error{e1}}
	return e.Append(e2)
}
