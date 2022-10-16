// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package array

import (
	"github.com/conflowio/conflow/src/schema"
)

// Contains returns true if the array contains the given element
// @function
func Contains(arr []interface{}, elem interface{}) (bool, error) {
	s, err := schema.GetSchemaForValue(arr)
	if err != nil {
		return false, err
	}

	itemSchema := s.(*schema.Array).Items
	for _, e := range arr {
		if itemSchema.CompareValues(e, elem) == 0 {
			return true, nil
		}
	}

	return false, nil
}
