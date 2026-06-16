// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json

import (
	"encoding/json"

	"github.com/conflowio/conflow/pkg/conflow/function"
	"github.com/conflowio/conflow/pkg/values"
)

// Encode converts the given value to a json string
// @function
func Encode(value interface{}) (string, error) {
	if slice, err := values.AsInterfaceSlice(value); err == nil {
		value = slice
	} else if goMap, err := values.AsStringInterfaceMap(value); err == nil {
		value = goMap
	}

	b, err := json.Marshal(value)
	if err != nil {
		return "", function.NewErrorf(0, "encoding JSON failed: %s", err)
	}
	return string(b), nil
}
