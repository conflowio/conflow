// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/conflowio/conflow/pkg/util/validation"
)

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type Any struct {
	Metadata

	Const    interface{}   `json:"const,omitempty"`
	Default  interface{}   `json:"default,omitempty"`
	Enum     []interface{} `json:"enum,omitempty"`
	Nullable bool          `json:"nullable,omitempty"`
}

func (a *Any) AssignValue(imports map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s", resultName, valueName)
}

func (a *Any) CompareValues(v1, v2 interface{}) int {
	panic("CompareValues should not be called on Any")
}

func (a *Any) Copy() Schema {
	j, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Any{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (a *Any) DefaultValue() interface{} {
	return a.Default
}

func (a *Any) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sAny{\n", schemaPkg(imports))
	if !reflect.ValueOf(a.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(a.Metadata.GoString(imports)))
	}
	if a.Const != nil {
		fprintf(buf, "\tConst: %#v,\n", a.Const)
	}
	if a.Default != nil {
		fprintf(buf, "\tDefault: %#v,\n", a.Default)
	}
	if len(a.Enum) > 0 {
		fprintf(buf, "\tEnum: %#v,\n", a.Enum)
	}
	if a.Nullable {
		fprintf(buf, "\tNullable: %#v,\n", a.Nullable)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (a *Any) GoType(_ map[string]string) string {
	return "interface{}"
}

func (a *Any) StringValue(value interface{}) string {
	s, err := GetSchemaForValue(value)
	if err != nil {
		return ""
	}
	return s.StringValue(value)
}

func (a *Any) Type() Type {
	return TypeAny
}

func (a *Any) TypeString() string {
	return "any"
}

func (a *Any) Validate(ctx context.Context) error {
	return validation.ValidateObject(ctx,
		validateCommonFields(a, a.Const, a.Default, a.Enum),
	)
}

func (a *Any) ValidateSchema(s Schema, compare bool) error {
	return nil
}

func (a *Any) ValidateValue(v interface{}) (interface{}, error) {
	if _, ok := v.(io.ReadCloser); ok {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, nil
		}
		v = rv.Elem().Interface()
	}

	if v == nil {
		return nil, nil
	}

	if a.Const != nil && a.Const != v {
		return nil, fmt.Errorf("must be %s", a.StringValue(a.Const))
	}

	if len(a.Enum) > 0 {
		if len(a.Enum) == 1 && a.Enum[0] != v {
			return nil, fmt.Errorf("must be %s", a.StringValue(a.Enum[0]))
		}

		allowed := func() bool {
			for _, e := range a.Enum {
				if e == v {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", a.join(a.Enum, ", "))
		}
	}

	return v, nil
}

func (a *Any) join(elems []interface{}, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return a.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(a.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(a.StringValue(e))
	}
	return b.String()
}

func AnyValue() Schema {
	return anyValueInst
}

var anyValueInst = &anyValue{
	Any: &Any{},
}

type anyValue struct {
	*Any
}

func (u *anyValue) Copy() Schema {
	return anyValueInst
}

func (u *anyValue) GoString(imports map[string]string) string {
	return fmt.Sprintf("%sAnyValue()", schemaPkg(imports))
}
