// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import "github.com/opsidian/basil/basil/variable"

// ArrayContains returns true if the array contains the given element
//go:generate basil generate
func ArrayContains(arr []interface{}, elem interface{}) bool {
	switch elem.(type) {
	case []interface{}, map[string]interface{}:
		for _, item := range arr {
			if variable.Equals(elem, item) {
				return true
			}
		}
	default:
		for _, item := range arr {
			if item == elem {
				return true
			}
		}
	}

	return false
}
