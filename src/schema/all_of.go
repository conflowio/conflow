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
)

type AllOf struct {
	// @name "schema"
	// @required
	// @min_items 1
	Schemas []Schema `json:"allOf"`
	// @ignore
	Schema `json:"-"`
}

func (a *AllOf) getSchema() Schema {
	if a.Schema == nil {
		panic("getSchema is not implemented yet")
	}
	return a.Schema
}

func (a *AllOf) AssignValue(imports map[string]string, valueName, resultName string) string {
	return a.getSchema().AssignValue(imports, valueName, resultName)
}

func (a *AllOf) CompareValues(v1, v2 interface{}) int {
	return a.getSchema().CompareValues(v1, v2)
}

func (a *AllOf) Copy() Schema {
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

func (a *AllOf) DefaultValue() interface{} {
	return a.getSchema().DefaultValue()
}

func (a *AllOf) GoType(imports map[string]string) string {
	return a.getSchema().GoType(imports)
}

func (a *AllOf) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.AllOf{\n")
	_, _ = fmt.Fprintf(buf, "\tSchemas: []schema.Schema{\n")
	for _, s := range a.Schemas {
		_, _ = fmt.Fprintf(buf, "\t\t%s,\n", indent(s.GoString(imports)))
	}
	_, _ = fmt.Fprintf(buf, "\t},\n")
	buf.WriteRune('}')
	return buf.String()
}

func (a *AllOf) StringValue(value interface{}) string {
	return a.getSchema().StringValue(value)
}

func (a *AllOf) Type() Type {
	return a.getSchema().Type()
}

func (a *AllOf) TypeString() string {
	return a.getSchema().TypeString()
}

func (a *AllOf) UnmarshalJSON(j []byte) error {
	allOf := struct {
		Schemas []*SchemaUnmarshaler `json:"allOf"`
	}{}
	if err := json.Unmarshal(j, &allOf); err != nil {
		return err
	}

	a.Schemas = make([]Schema, len(allOf.Schemas))
	for _, s := range allOf.Schemas {
		a.Schemas = append(a.Schemas, s.Schema)
	}

	return nil
}

func (a *AllOf) ValidateSchema(s Schema, compare bool) error {
	return a.getSchema().ValidateSchema(s, compare)
}

func (a *AllOf) ValidateValue(value interface{}) (interface{}, error) {
	return a.getSchema().ValidateValue(value)
}
