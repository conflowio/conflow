// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bind

import (
	"fmt"
	"os"
	"reflect"

	"github.com/conflowio/conflow/pkg/schema"
)

var bindDebug = os.Getenv("CONFLOW_BIND_DEBUG") != ""

func debugBind(schemaType schema.Type, value interface{}) {
	if !bindDebug {
		return
	}
	fmt.Fprintf(os.Stderr, "conflow bind: schema=%s value=%s\n", schemaType, valueKind(value))
}

func valueKind(value interface{}) string {
	if value == nil {
		return "nil"
	}
	if isValuesList(value) {
		return "values.List"
	}
	if isValuesMap(value) {
		return "values.Map"
	}
	if frozen, ok := freezeListBuilder(value); ok {
		return "values.ListBuilder -> " + valueKind(frozen)
	}
	if frozen, ok := freezeMapBuilder(value); ok {
		return "values.MapBuilder -> " + valueKind(frozen)
	}
	return reflect.TypeOf(value).String()
}
