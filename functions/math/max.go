// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package math

import (
	"github.com/conflowio/conflow/conflow/schema"
)

// Max returns with the greatest value
// @function
func Max(
	// @types ["integer", "number"]
	max interface{},
	// @types ["integer", "number"]
	rest ...interface{},
) float64 {
	var res float64
	switch maxt := max.(type) {
	case int64:
		res = float64(maxt)
	case float64:
		res = maxt
	}

	for _, v := range rest {
		switch vt := v.(type) {
		case int64:
			if schema.NumberGreaterThan(float64(vt), res) {
				res = float64(vt)
			}
		case float64:
			if schema.NumberGreaterThan(vt, res) {
				res = vt
			}
		}
	}
	return res
}
