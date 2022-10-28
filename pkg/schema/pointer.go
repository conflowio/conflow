// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"fmt"
	"strings"
)

func Pointer[V any](v V) *V {
	return &v
}

func PointerValue[V any](v interface{}) *V {
	switch vt := v.(type) {
	case *V:
		return vt
	case V:
		return &vt
	default:
		var n V
		panic(fmt.Errorf("unexpected value %T, was expecting %T", v, n))
	}
}

func Value[V any](v interface{}) (r V) {
	switch vt := v.(type) {
	case *V:
		if vt != nil {
			return *vt
		}
		return r
	case V:
		return vt
	default:
		panic(fmt.Errorf("unexpected value %T, was expecting %T", v, r))
	}
}

func assignFuncName(schema Schema, imports map[string]string) string {
	t := strings.TrimPrefix(schema.GoType(imports), "*")
	if n, ok := schema.(Nullable); ok && n.GetNullable() {
		return fmt.Sprintf("%sPointerValue[%s]", schemaPkg(imports), t)
	}

	return fmt.Sprintf("%sValue[%s]", schemaPkg(imports), t)
}
