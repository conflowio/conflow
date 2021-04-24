// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Integer struct {
	Metadata

	Const   *int64  `json:"const,omitempty"`
	Default *int64  `json:"default,omitempty"`
	Enum    []int64 `json:"enum,omitempty"`
}

func (i *Integer) AssignValue(imports map[string]string, valueName, resultName string) string {
	if i.Pointer {
		schemaPackageName := EnsureUniqueGoPackageName(imports, "github.com/opsidian/basil/basil/schema")
		return fmt.Sprintf("%s = %s.IntegerPtr(%s.(int64))", resultName, schemaPackageName, valueName)
	}

	return fmt.Sprintf("%s = %s.(int64)", resultName, valueName)
}

func (i *Integer) CompareValues(v1, v2 interface{}) int {
	var f1 float64
	switch v := v1.(type) {
	case int64:
		f1 = float64(v)
	case float64:
		f1 = v
	default:
		return -1
	}

	var f2 float64
	switch v := v2.(type) {
	case int64:
		f2 = float64(v)
	case float64:
		f2 = v
	default:
		return -1
	}

	switch {
	case f1-f2 < Epsilon && f2-f1 < Epsilon:
		return 0
	case f1 < f2:
		return -1
	default:
		return 1
	}
}

func (i *Integer) Copy() Schema {
	j, err := json.Marshal(i)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Integer{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (i *Integer) DefaultValue() interface{} {
	if i.Default == nil {
		return nil
	}
	return *i.Default
}

func (i *Integer) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Integer{\n")
	if !reflect.ValueOf(i.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(i.Metadata.GoString()))
	}
	if i.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: schema.IntegerPtr(%#v),\n", *i.Const)
	}
	if i.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: schema.IntegerPtr(%#v),\n", *i.Default)
	}
	if len(i.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", i.Enum)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (i *Integer) GoType(_ map[string]string) string {
	if i.Pointer {
		return "*int64"
	}
	return "int64"
}

func (i *Integer) MarshalJSON() ([]byte, error) {
	type Alias Integer
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(i.Type()),
		Alias: (*Alias)(i),
	})
}

func (i *Integer) StringValue(value interface{}) string {
	v, ok := value.(int64)
	if !ok {
		return ""
	}

	return strconv.FormatInt(v, 10)
}

func (i *Integer) Type() Type {
	return TypeInteger
}

func (i *Integer) TypeString() string {
	return string(TypeInteger)
}

func (i *Integer) ValidateSchema(i2 Schema, compare bool) error {
	if compare {
		if i2.Type() != TypeInteger && i2.Type() != TypeNumber {
			return typeError("must be number")
		}
		return nil
	}

	if i2.Type() != TypeInteger {
		return typeError("must be integer")
	}

	return nil
}

func (i *Integer) ValidateValue(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("must be integer")
	}

	if i.Const != nil && *i.Const != v {
		return fmt.Errorf("must be %s", i.StringValue(*i.Const))
	}

	if len(i.Enum) == 1 && i.Enum[0] != v {
		return fmt.Errorf("must be %s", i.StringValue(i.Enum[0]))
	}

	if len(i.Enum) > 0 {
		allowed := func() bool {
			for _, e := range i.Enum {
				if e == v {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return fmt.Errorf("must be one of %s", i.join(i.Enum, ", "))
		}
	}

	return nil
}

func (i *Integer) join(elems []int64, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return i.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(i.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(i.StringValue(e))
	}
	return b.String()
}

func IntegerValue() Schema {
	return integerValueInst
}

var integerValueInst = &integerValue{
	Integer: &Integer{},
}

type integerValue struct {
	*Integer
}

func (i *integerValue) Copy() Schema {
	return integerValueInst
}

func (i *integerValue) GoString() string {
	return "schema.IntegerValue()"
}

func IntegerPtr(v int64) *int64 {
	return &v
}
