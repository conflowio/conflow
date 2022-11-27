// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions

import (
	"fmt"
	"unicode/utf8"
)

// Len returns with the length of the variable
// For strings it means the count of UTF-8 characters
// For arrays and maps it means the number of items/entries
// @function
func Len(
	// @one_of {
	//   schema:string
	//   schema:array {
	//     items:any
	//   }
	//   schema:map {
	//     additional_properties:any
	//   }
	// }
	value interface{},
) int64 {
	switch v := value.(type) {
	case string:
		return int64(utf8.RuneCountInString(v))
	case []interface{}:
		return int64(len(v))
	case map[string]interface{}:
		return int64(len(v))
	default:
		panic(fmt.Sprintf("unexpected type: %T", v))
	}
}
