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

// Epsilon is used as a float64 comparison tolerance
const Epsilon = 0.000000001

type Number struct {
	Metadata

	Const   *float64  `json:"const,omitempty"`
	Default *float64  `json:"default,omitempty"`
	Enum    []float64 `json:"enum,omitempty"`
}

func (n *Number) AssignValue(imports map[string]string, valueName, resultName string) string {
	if n.Pointer {
		schemaPackageName := EnsureUniqueGoPackageName(imports, "github.com/opsidian/basil/basil/schema")
		return fmt.Sprintf("%s = %s.NumberPtr(%s.(float64))", resultName, schemaPackageName, valueName)
	}

	return fmt.Sprintf("%s = %s.(float64)", resultName, valueName)
}

func (n *Number) CompareValues(v1, v2 interface{}) int {
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

func (n *Number) Copy() Schema {
	j, err := json.Marshal(n)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Number{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (n *Number) DefaultValue() interface{} {
	if n.Default == nil {
		return nil
	}
	return *n.Default
}

func (n *Number) MarshalJSON() ([]byte, error) {
	type Alias Number
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(n.Type()),
		Alias: (*Alias)(n),
	})
}

func (n *Number) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Number{\n")
	if !reflect.ValueOf(n.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(n.Metadata.GoString()))
	}
	if n.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: schema.NumberPtr(%#v),\n", *n.Const)
	}
	if n.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: schema.NumberPtr(%#v),\n", *n.Default)
	}
	if len(n.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", n.Enum)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (n *Number) GoType(_ map[string]string) string {
	if n.Pointer {
		return "*float64"
	}
	return "float64"
}

func (n *Number) StringValue(value interface{}) string {
	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func (n *Number) Type() Type {
	return TypeNumber
}

func (n *Number) TypeString() string {
	return string(TypeNumber)
}

func (n *Number) ValidateSchema(n2 Schema, _ bool) error {
	if n2.Type() != TypeNumber && n2.Type() != TypeInteger {
		return typeError("must be number")
	}

	return nil
}

func (n *Number) ValidateValue(value interface{}) error {
	var v float64
	switch vt := value.(type) {
	case int64:
		v = float64(vt)
	case float64:
		v = vt
	default:
		return errors.New("must be number")
	}

	if n.Const != nil && *n.Const != v {
		return fmt.Errorf("must be %s", n.StringValue(*n.Const))
	}

	if len(n.Enum) == 1 && n.Enum[0] != v {
		return fmt.Errorf("must be %s", n.StringValue(n.Enum[0]))
	}

	if len(n.Enum) > 0 {
		allowed := func() bool {
			for _, e := range n.Enum {
				if e == v {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return fmt.Errorf("must be one of %s", n.join(n.Enum, ", "))
		}
	}

	return nil
}

func (n *Number) join(elems []float64, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return n.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(n.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(n.StringValue(e))
	}
	return b.String()
}

func NumberValue() Schema {
	return numberValueInst
}

var numberValueInst = &numberValue{
	Number: &Number{},
}

type numberValue struct {
	*Number
}

func (n *numberValue) Copy() Schema {
	return numberValueInst
}

func (n *numberValue) GoString() string {
	return "schema.NumberValue()"
}

func NumberPtr(v float64) *float64 {
	return &v
}
