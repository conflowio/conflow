// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package math

import (
	"fmt"
)

// Abs returns the absolute value of the given number
// @function
func Abs(
	// @one_of {
	//   schema:integer
	//   schema:number
	// }
	// @result_type
	value interface{},
) interface{} {
	switch n := value.(type) {
	case int64:
		if n >= 0 {
			return n
		}
		return -1 * n
	case float64:
		if n >= 0 {
			return value
		}
		return -1 * n
	default:
		panic(fmt.Sprintf("unexpected value type: %T", value))
	}
}
