// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions

import (
	"github.com/conflowio/conflow/conflow/schema"
)

// ArrayContains returns true if the array contains the given element
// @function
func ArrayContains(arr []interface{}, elem interface{}) (bool, error) {
	s, err := schema.GetSchemaForValue(arr)
	if err != nil {
		return false, err
	}

	itemSchema := s.(schema.ArrayKind).GetItems()
	for _, e := range arr {
		if itemSchema.CompareValues(e, elem) == 0 {
			return true, nil
		}
	}

	return false, nil
}
