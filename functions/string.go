// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package functions

import (
	"fmt"
	"strconv"

	"github.com/opsidian/basil/basil/variable"
)

// String converts the given value to a string
//go:generate basil generate
func String(value variable.Basic) string {
	switch v := value.Value().(type) {
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		panic(fmt.Sprintf("unexpected value type: %T", value.Value()))
	}
}
