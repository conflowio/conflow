// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package math

import (
	"fmt"
	"math"

	"github.com/opsidian/basil/basil/variable"
)

// Floor returns the greatest integer value less than or equal to x.
//go:generate basil generate
func Floor(number *variable.Number) int64 {
	switch v := number.Value().(type) {
	case int64:
		return v
	case float64:
		return int64(math.Floor(v))
	default:
		panic(fmt.Sprintf("unexpected type: %T", number.Value()))
	}
}
