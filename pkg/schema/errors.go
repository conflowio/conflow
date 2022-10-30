// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"fmt"
)

type typeError string

func (t typeError) Error() string {
	return string(t)
}

func typeErrorf(format string, a ...interface{}) typeError {
	return typeError(fmt.Sprintf(format, a...))
}
