// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package math

import (
	"math"
)

// Round returns the nearest integer, rounding half away from zero.
// @function
func Round(value float64) float64 {
	return math.Round(value)
}
