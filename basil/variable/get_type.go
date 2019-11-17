// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package variable

import (
	"io"
	"time"

	"github.com/opsidian/basil/basil"
)

// GetType returns with the type of the given value
func GetType(value interface{}) string {
	switch value.(type) {
	case []interface{}:
		return TypeArray
	case *Basic:
		return TypeBasic
	case bool:
		return TypeBool
	case float64:
		return TypeFloat
	case basil.ID:
		return TypeIdentifier
	case int64:
		return TypeInteger
	case map[string]interface{}:
		return TypeMap
	case *Number:
		return TypeNumber
	case io.ReadCloser:
		return TypeStream
	case string:
		return TypeString
	case []string:
		return TypeStringArray
	case time.Time:
		return TypeTime
	case time.Duration:
		return TypeTimeDuration
	case *WithLength:
		return TypeWithLength
	default:
		return TypeUnknown
	}
}
