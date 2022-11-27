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
	"reflect"
	"sort"
	"strings"

	"github.com/conflowio/conflow/pkg/util/validation"
)

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type AnyOf struct {
	Metadata

	Const    interface{}   `json:"const,omitempty"`
	Default  interface{}   `json:"default,omitempty"`
	Enum     []interface{} `json:"enum,omitempty"`
	Nullable bool          `json:"nullable,omitempty"`

	// @name "schema"
	// @required
	// @min_items 1
	Schemas []Schema `json:"anyOf"`
}

func (a *AnyOf) AssignValue(_ map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s", resultName, valueName)
}

func (a *AnyOf) CompareValues(_, _ interface{}) int {
	panic("CompareValues should not be called on AnyOf")
}

func (a *AnyOf) Copy() Schema {
	j, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &AnyOf{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (a *AnyOf) DefaultValue() interface{} {
	return a.Default
}

func (a *AnyOf) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sAnyOf{\n", schemaPkg(imports))
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
	if len(a.Schemas) > 0 {
		fprintf(buf, "\tSchemas: %s,\n", indent(a.schemasString(imports)))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (a *AnyOf) GoType(imports map[string]string) string {
	if len(a.Schemas) == 1 {
		return a.Schemas[0].GoType(imports)
	}
	return "interface{}"
}

func (a *AnyOf) GetNullable() bool {
	return a.Nullable
}

func (a *AnyOf) SetNullable(nullable bool) {
	a.Nullable = nullable
}

func (a *AnyOf) StringValue(value interface{}) string {
	for _, s := range a.Schemas {
		if vv, err := s.ValidateValue(value); err == nil {
			return s.StringValue(vv)
		}
	}

	return ""
}

func (a *AnyOf) Type() Type {
	if len(a.Schemas) == 1 {
		return a.Schemas[0].Type()
	}

	return TypeAny
}

func (a *AnyOf) TypeString() string {
	if len(a.Schemas) == 1 {
		return a.Schemas[0].TypeString()
	}

	var types []string
	for _, s := range a.Schemas {
		types = append(types, s.TypeString())
	}
	sort.Strings(types)
	return fmt.Sprintf(
		"%s or %s",
		strings.Join(types[0:len(types)-1], ", "),
		types[len(types)-1],
	)
}

func (a *AnyOf) UnmarshalJSON(j []byte) error {
	type Alias AnyOf
	v := struct {
		*Alias
		Schemas []*SchemaUnmarshaler `json:"anyOf"`
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

func (a *AnyOf) Validate(ctx context.Context) error {
	return validation.ValidateObject(ctx,
		validation.ValidateArray("", a.Schemas),
		validateCommonFields(a, a.Const, a.Default, a.Enum),
	)
}

func (a *AnyOf) ValidateSchema(schema Schema, compare bool) error {
	for _, s := range a.Schemas {
		if err := s.ValidateSchema(schema, compare); err == nil {
			return nil
		}
	}

	return typeErrorf("was expecting %s", a.TypeString())
}

func (a *AnyOf) ValidateValue(v interface{}) (interface{}, error) {
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

	for _, s := range a.Schemas {
		vv, err := s.ValidateValue(v)
		if err == nil {
			return vv, nil
		}
	}

	return nil, fmt.Errorf("must be %s", a.TypeString())
}

func (a *AnyOf) schemasString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "[]%sSchema{\n", schemaPkg(imports))
	for _, s := range a.Schemas {
		fprintf(buf, "\t%s,\n", indent(s.GoString(imports)))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (a *AnyOf) join(elems []interface{}, sep string) string {
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
