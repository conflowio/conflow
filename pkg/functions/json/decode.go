// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package json

import (
	"encoding/json"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow/function"
)

// Decode converts the given json string to a data structure
// @function
func Decode(jsonStr string) (interface{}, error) {
	var val interface{}
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	dec.UseNumber()
	if err := dec.Decode(&val); err != nil {
		return nil, function.NewErrorf(0, "decoding JSON failed: %s", err)
	}
	return convertJSONNumbers(val), nil
}

func convertJSONNumbers(val interface{}) interface{} {
	switch v := val.(type) {
	case json.Number:
		if intVal, err := v.Int64(); err == nil {
			return intVal
		}
		floatVal, _ := v.Float64()
		return floatVal
	case []interface{}:
		for vk, vv := range v {
			v[vk] = convertJSONNumbers(vv)
		}
	case map[string]interface{}:
		for vk, vv := range v {
			v[vk] = convertJSONNumbers(vv)
		}
	}
	return val
}
