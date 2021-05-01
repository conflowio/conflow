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

type Array struct {
	Metadata

	Const   []interface{}   `json:"const,omitempty"`
	Default []interface{}   `json:"default,omitempty"`
	Enum    [][]interface{} `json:"enum,omitempty"`
	// @required
	Items Schema `json:"items,omitempty"`
}

func (a *Array) AssignValue(imports map[string]string, valueName, resultName string) string {
	if a.Pointer {
		panic("an array value can not have a pointer")
	}

	if a.Items.Type() == TypeUntyped {
		return fmt.Sprintf("%s = %s.([]interface{})", resultName, valueName)
	}

	return fmt.Sprintf(`%s = make(%s, len(%s.([]interface{})))
for %sk, %sv := range %s.([]interface{}) {
	%s
}`,
		resultName,
		a.GoType(imports),
		valueName,
		valueName,
		valueName,
		valueName,
		indent(a.Items.AssignValue(imports, valueName+"v", fmt.Sprintf("%s[%sk]", resultName, valueName))),
	)
}

func (a *Array) CompareValues(v1, v2 interface{}) int {
	var a1 []interface{}
	if v1 != nil {
		var ok bool
		if a1, ok = v1.([]interface{}); !ok {
			return -1
		}
	}

	var a2 []interface{}
	if v2 != nil {
		var ok bool
		if a2, ok = v2.([]interface{}); !ok {
			return 1
		}
	}

	switch {
	case len(a1) == len(a2):
		for i := 0; i < len(a1); i++ {
			if c := a.Items.CompareValues(a1[i], a2[i]); c != 0 {
				return c
			}
		}
		return 0
	case len(a1) < len(a2):
		return -1
	default:
		return 1
	}
}

func (a *Array) Copy() Schema {
	j, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Array{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (a *Array) DefaultValue() interface{} {
	if a.Default == nil {
		return nil
	}

	return a.Default
}

func (a *Array) GetItems() Schema {
	return a.Items
}

func (a *Array) GoType(imports map[string]string) string {
	if a.Pointer {
		panic("an array value can not have a pointer")
	}

	return fmt.Sprintf("[]%s", a.Items.GoType(imports))
}

func (a *Array) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Array{\n")
	if !reflect.ValueOf(a.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(a.Metadata.GoString()))
	}
	if a.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: %#v,\n", a.Const)
	}
	if a.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: %#v,\n", a.Default)
	}
	if len(a.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", a.Enum)
	}
	if a.Items != nil {
		_, _ = fmt.Fprintf(buf, "\tItems: %s,\n", indent(a.Items.GoString()))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (a *Array) MarshalJSON() ([]byte, error) {
	type Alias Array
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(a.Type()),
		Alias: (*Alias)(a),
	})
}

func (a *Array) StringValue(value interface{}) string {
	v, ok := value.([]interface{})
	if !ok {
		return ""
	}

	if len(v) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteRune('[')
	sb.WriteString(a.Items.StringValue(v[0]))
	for _, e := range v[1:] {
		sb.WriteString(", ")
		sb.WriteString(a.Items.StringValue(e))
	}
	sb.WriteRune(']')
	return sb.String()
}

func (a *Array) Type() Type {
	return TypeArray
}

func (a *Array) TypeString() string {
	if a.Items.Type() == TypeUntyped {
		return "array"
	}

	return fmt.Sprintf("array(%s)", a.Items.TypeString())
}

func (a *Array) UnmarshalJSON(j []byte) error {
	type Alias Array
	v := struct {
		*Alias
		Items *SchemaUnmarshaler `json:"items,omitempty"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	a.Items = v.Items.Schema

	return nil
}

func (a *Array) ValidateSchema(s Schema, compare bool) error {
	if s.Type() == TypeNull {
		return nil
	}

	a2, ok := s.(ArrayKind)
	if !ok {
		return typeErrorf("must be %s", a.TypeString())
	}

	if err := a.Items.ValidateSchema(a2.GetItems(), compare); err != nil {
		if _, ok := err.(typeError); ok {
			return typeErrorf("must be %s", a.TypeString())
		}
		return err
	}

	return nil
}

func (a *Array) ValidateValue(value interface{}) error {
	var v []interface{}
	if value != nil {
		var ok bool
		if v, ok = value.([]interface{}); !ok {
			return errors.New("must be array")
		}
	}

	if a.Const != nil && a.CompareValues(a.Const, v) != 0 {
		return fmt.Errorf("must be %s", a.StringValue(a.Const))
	}

	if len(a.Enum) == 1 && a.CompareValues(a.Enum[0], v) != 0 {
		return fmt.Errorf("must be %s", a.StringValue(a.Enum[0]))
	}

	if len(a.Enum) > 0 {
		allowed := func() bool {
			for _, e := range a.Enum {
				if a.CompareValues(e, v) == 0 {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return fmt.Errorf("must be one of %s", a.join(a.Enum, ", "))
		}
	}

	ve := ValidationError{}
	for i, e := range v {
		if err := a.Items.ValidateValue(e); err != nil {
			ve.AddError(strconv.Itoa(i), err)
		}
	}

	return ve.ErrOrNil()
}

func (a *Array) join(elems [][]interface{}, sep string) string {
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
