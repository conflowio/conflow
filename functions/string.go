// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions

import (
	"github.com/opsidian/conflow/conflow/schema"
)

// String converts the given value to a string
// @function
func String(value interface{}) (string, error) {
	s, err := schema.GetSchemaForValue(value)
	if err != nil {
		return "", err
	}

	return s.StringValue(value), nil
}
