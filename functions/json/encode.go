// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json

import (
	"encoding/json"

	"github.com/opsidian/conflow/conflow/function"
)

// Encode converts the given value to a json string
// @function
func Encode(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", function.NewErrorf(0, "encoding JSON failed: %s", err)
	}
	return string(b), nil
}
