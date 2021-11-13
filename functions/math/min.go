// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package math

import (
	"github.com/conflowio/conflow/conflow/schema"
)

// Min returns with the lowest value
// @function
func Min(
	// @types ["integer", "number"]
	min interface{},
	// @types ["integer", "number"]
	rest ...interface{},
) float64 {
	var res float64
	switch mint := min.(type) {
	case int64:
		res = float64(mint)
	case float64:
		res = mint
	}

	for _, v := range rest {
		switch vt := v.(type) {
		case int64:
			if schema.NumberLessThan(float64(vt), res) {
				res = float64(vt)
			}
		case float64:
			if schema.NumberLessThan(vt, res) {
				res = vt
			}
		}
	}
	return res
}
