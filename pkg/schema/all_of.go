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

	"github.com/conflowio/conflow/pkg/util/validation"
)

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type AllOf struct {
	Const    interface{}   `json:"const,omitempty"`
	Default  interface{}   `json:"default,omitempty"`
	Enum     []interface{} `json:"enum,omitempty"`
	Nullable bool          `json:"nullable,omitempty"`

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

	cp := &AllOf{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (a *AllOf) DefaultValue() interface{} {
	return a.Default
}

func (a *AllOf) GoType(imports map[string]string) string {
	return a.getSchema().GoType(imports)
}

func (a *AllOf) GoString(imports map[string]string) string {
	pkg := schemaPkg(imports)
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sAllOf{\n", pkg)
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
	fprintf(buf, "\tSchemas: []%sSchema{\n", pkg)
	for _, s := range a.Schemas {
		fprintf(buf, "\t\t%s,\n", indent(s.GoString(imports)))
	}
	fprintf(buf, "\t},\n")
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
	type Alias AllOf
	v := struct {
		*Alias
		Schemas []*SchemaUnmarshaler `json:"allOf"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	a.Schemas = make([]Schema, 0, len(v.Schemas))
	for _, s := range v.Schemas {
		a.Schemas = append(a.Schemas, s.Schema)
	}

	return nil
}

func (a *AllOf) Validate(ctx context.Context) error {
	return validation.ValidateObject(ctx,
		validation.ValidateArray("", a.Schemas),
	)
}

func (a *AllOf) ValidateSchema(s Schema, compare bool) error {
	return a.getSchema().ValidateSchema(s, compare)
}

func (a *AllOf) ValidateValue(value interface{}) (interface{}, error) {
	return a.getSchema().ValidateValue(value)
}
