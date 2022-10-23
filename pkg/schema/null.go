// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type Null struct {
	Metadata
}

func (n *Null) AssignValue(_ map[string]string, _, resultName string) string {
	return fmt.Sprintf("%s = nil", resultName)
}

func (n *Null) CompareValues(v1, v2 interface{}) int {
	if v1 == nil && v2 == nil {
		return 0
	}

	return -1
}

func (n *Null) Copy() Schema {
	j, err := json.Marshal(n)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Null{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (n *Null) DefaultValue() interface{} {
	return nil
}

func (n *Null) MarshalJSON() ([]byte, error) {
	type Alias Null
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(n.Type()),
		Alias: (*Alias)(n),
	})
}

func (n *Null) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sNull{\n", schemaPkg(imports))
	if !reflect.ValueOf(n.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(n.Metadata.GoString(imports)))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (n *Null) GoType(_ map[string]string) string {
	return "nil"
}

func (n *Null) StringValue(interface{}) string {
	return "null"
}

func (n *Null) Type() Type {
	return TypeNull
}

func (n *Null) TypeString() string {
	return string(TypeNull)
}

func (n *Null) Validate(*Context) error {
	return nil
}

func (n *Null) ValidateSchema(n2 Schema, _ bool) error {
	if n2.Type() != TypeNull && n2.Type() != TypeArray && n2.Type() != TypeMap {
		return typeError("must be null, array or map")
	}

	return nil
}

func (n *Null) ValidateValue(v interface{}) (interface{}, error) {
	switch vt := v.(type) {
	case nil:
		return nil, nil
	case []interface{}:
		if len(vt) == 0 {
			return []interface{}{}, nil
		}
	case map[string]interface{}:
		if len(vt) == 0 {
			return map[string]interface{}{}, nil
		}
	}

	return nil, fmt.Errorf("must be null, empty array or empty map")
}

func NullValue() Schema {
	return nullValue
}

var nullValue = &struct {
	*Null
}{
	Null: &Null{},
}
