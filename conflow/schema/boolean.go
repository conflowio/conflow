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
)

type Boolean struct {
	Metadata

	Const   *bool  `json:"const,omitempty"`
	Default *bool  `json:"default,omitempty"`
	Enum    []bool `json:"enum,omitempty"`
}

func (b *Boolean) AssignValue(imports map[string]string, valueName, resultName string) string {
	if b.Pointer {
		schemaPackageName := EnsureUniqueGoPackageName(imports, "github.com/conflowio/conflow/conflow/schema")
		return fmt.Sprintf("%s = %s.BooleanPtr(%s.(bool))", resultName, schemaPackageName, valueName)
	}

	return fmt.Sprintf("%s = %s.(bool)", resultName, valueName)
}

func (b *Boolean) CompareValues(v1, v2 interface{}) int {
	b1, ok := v1.(bool)
	if !ok {
		return -1
	}

	b2, ok := v2.(bool)
	if !ok {
		return 1
	}

	switch {
	case b1 == b2:
		return 0
	case !b1:
		return -1
	default:
		return 1
	}
}

func (b *Boolean) Copy() Schema {
	j, err := json.Marshal(b)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Boolean{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (b *Boolean) DefaultValue() interface{} {
	if b.Default == nil {
		return nil
	}
	return *b.Default
}

func (b *Boolean) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Boolean{\n")
	if !reflect.ValueOf(b.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(b.Metadata.GoString()))
	}
	if b.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: schema.BooleanPtr(%#v),\n", *b.Const)
	}
	if b.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: schema.BooleanPtr(%#v),\n", *b.Default)
	}
	if len(b.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", b.Enum)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (b *Boolean) GoType(_ map[string]string) string {
	if b.Pointer {
		return "*bool"
	}
	return "bool"
}

func (b *Boolean) MarshalJSON() ([]byte, error) {
	type Alias Boolean
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(b.Type()),
		Alias: (*Alias)(b),
	})
}

func (b *Boolean) StringValue(value interface{}) string {
	v, ok := value.(bool)
	if !ok {
		return ""
	}

	return strconv.FormatBool(v)
}

func (b *Boolean) Type() Type {
	return TypeBoolean
}

func (b *Boolean) TypeString() string {
	return string(TypeBoolean)
}

func (b *Boolean) ValidateSchema(b2 Schema, _ bool) error {
	if b2.Type() != TypeBoolean {
		return typeError("must be boolean")
	}

	return nil
}

func (b *Boolean) ValidateValue(value interface{}) error {
	v, ok := value.(bool)
	if !ok {
		return errors.New("must be boolean")
	}

	if b.Const != nil && *b.Const != v {
		return fmt.Errorf("must be %s", b.StringValue(*b.Const))
	}

	if len(b.Enum) == 1 && b.Enum[0] != v {
		return fmt.Errorf("must be %s", b.StringValue(b.Enum[0]))
	}

	return nil
}

func BooleanValue() Schema {
	return booleanValueInst
}

var booleanValueInst = &booleanValue{
	Boolean: &Boolean{},
}

type booleanValue struct {
	*Boolean
}

func (b *booleanValue) Copy() Schema {
	return booleanValueInst
}

func (b *booleanValue) GoString() string {
	return "schema.BooleanValue()"
}

func BooleanPtr(v bool) *bool {
	return &v
}
