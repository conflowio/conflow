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
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type Boolean struct {
	Metadata

	Const    *bool  `json:"const,omitempty"`
	Default  *bool  `json:"default,omitempty"`
	Enum     []bool `json:"enum,omitempty"`
	Nullable bool   `json:"nullable,omitempty"`
}

func (b *Boolean) AssignValue(imports map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s(%s)", resultName, assignFuncName(b, imports), valueName)
}

func (b *Boolean) CompareValues(v1, v2 interface{}) int {
	b1, _ := b.valueOf(v1)
	b2, _ := b.valueOf(v2)

	switch {
	case b1 == nil && b2 == nil:
		return 0
	case b1 == nil:
		return -1
	case b2 == nil:
		return 1
	case *b1 == *b2:
		return 0
	case !*b1:
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

func (b *Boolean) GetNullable() bool {
	return b.Nullable
}

func (b *Boolean) GoString(imports map[string]string) string {
	pkg := schemaPkg(imports)
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sBoolean{\n", pkg)
	if !reflect.ValueOf(b.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(b.Metadata.GoString(imports)))
	}
	if b.Const != nil {
		fprintf(buf, "\tConst: %sPointer(%#v),\n", pkg, *b.Const)
	}
	if b.Default != nil {
		fprintf(buf, "\tDefault: %sPointer(%#v),\n", pkg, *b.Default)
	}
	if len(b.Enum) > 0 {
		fprintf(buf, "\tEnum: %#v,\n", b.Enum)
	}
	if b.Nullable {
		fprintf(buf, "\tNullable: %#v,\n", b.Nullable)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (b *Boolean) GoType(_ map[string]string) string {
	if b.Nullable {
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

func (b *Boolean) SetNullable(nullable bool) {
	b.Nullable = nullable
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

func (b *Boolean) UnmarshalJSON(input []byte) error {
	type Alias Boolean
	return json.Unmarshal(input, &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(b),
	})
}

func (b *Boolean) Validate(ctx context.Context) error {
	if err := validateCommonFields(b, b.Const, b.Default, b.Enum); err != nil {
		return err
	}

	return nil
}

func (b *Boolean) ValidateSchema(b2 Schema, _ bool) error {
	if b2.Type() != TypeBoolean {
		return typeError("must be boolean")
	}

	return nil
}

func (b *Boolean) ValidateValue(value interface{}) (interface{}, error) {
	v, ok := b.valueOf(value)
	if !ok {
		return nil, errors.New("must be boolean")
	}
	if v == nil {
		return nil, nil
	}

	if b.Const != nil && *b.Const != *v {
		return nil, fmt.Errorf("must be %s", b.StringValue(*b.Const))
	}

	if len(b.Enum) == 1 && b.Enum[0] != *v {
		return nil, fmt.Errorf("must be %s", b.StringValue(b.Enum[0]))
	}

	if b.Nullable {
		return v, nil
	}
	return *v, nil
}

func (b *Boolean) valueOf(value interface{}) (*bool, bool) {
	switch v := value.(type) {
	case bool:
		return &v, true
	case *bool:
		return v, true
	default:
		return nil, false
	}
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

func (b *booleanValue) GoString(imports map[string]string) string {
	return fmt.Sprintf("%sBooleanValue()", schemaPkg(imports))
}
