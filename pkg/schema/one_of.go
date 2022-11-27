// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozillo.org/MPL/2.0/.

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
type OneOf struct {
	Metadata

	Const    interface{}   `json:"const,omitempty"`
	Default  interface{}   `json:"default,omitempty"`
	Enum     []interface{} `json:"enum,omitempty"`
	Nullable bool          `json:"nullable,omitempty"`

	// @name "schema"
	// @required
	// @min_items 1
	Schemas []Schema `json:"oneOf"`
}

func (o *OneOf) AssignValue(_ map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s", resultName, valueName)
}

func (o *OneOf) CompareValues(_, _ interface{}) int {
	panic("CompareValues should not be called on OneOf")
}

func (o *OneOf) Copy() Schema {
	j, err := json.Marshal(o)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &OneOf{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (o *OneOf) DefaultValue() interface{} {
	return o.Default
}

func (o *OneOf) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sOneOf{\n", schemaPkg(imports))
	if !reflect.ValueOf(o.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(o.Metadata.GoString(imports)))
	}
	if o.Const != nil {
		fprintf(buf, "\tConst: %#v,\n", o.Const)
	}
	if o.Default != nil {
		fprintf(buf, "\tDefault: %#v,\n", o.Default)
	}
	if len(o.Enum) > 0 {
		fprintf(buf, "\tEnum: %#v,\n", o.Enum)
	}
	if o.Nullable {
		fprintf(buf, "\tNullable: %#v,\n", o.Nullable)
	}
	if len(o.Schemas) > 0 {
		fprintf(buf, "\tSchemas: %s,\n", indent(o.schemasString(imports)))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (o *OneOf) GoType(imports map[string]string) string {
	if len(o.Schemas) == 1 {
		return o.Schemas[0].GoType(imports)
	}
	return "interface{}"
}

func (o *OneOf) GetNullable() bool {
	return o.Nullable
}

func (o *OneOf) SetNullable(nullable bool) {
	o.Nullable = nullable
}

func (o *OneOf) StringValue(value interface{}) string {
	for _, s := range o.Schemas {
		if vv, err := s.ValidateValue(value); err == nil {
			return s.StringValue(vv)
		}
	}

	return ""
}

func (o *OneOf) Type() Type {
	if len(o.Schemas) == 1 {
		return o.Schemas[0].Type()
	}

	return TypeAny
}

func (o *OneOf) TypeString() string {
	if len(o.Schemas) == 1 {
		return o.Schemas[0].TypeString()
	}

	var types []string
	for _, s := range o.Schemas {
		types = append(types, s.TypeString())
	}
	sort.Strings(types)
	return fmt.Sprintf(
		"%s or %s",
		strings.Join(types[0:len(types)-1], ", "),
		types[len(types)-1],
	)
}

func (o *OneOf) UnmarshalJSON(j []byte) error {
	type Alias OneOf
	v := struct {
		*Alias
		Schemas []*SchemaUnmarshaler `json:"oneOf"`
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	o.Schemas = make([]Schema, 0, len(v.Schemas))
	for _, s := range v.Schemas {
		o.Schemas = append(o.Schemas, s.Schema)
	}

	return nil
}

func (o *OneOf) Validate(ctx context.Context) error {
	return validation.ValidateObject(ctx,
		validation.ValidateArray("", o.Schemas),
		validateCommonFields(o, o.Const, o.Default, o.Enum),
	)
}

func (o *OneOf) ValidateSchema(schema Schema, compare bool) error {
	for _, s := range o.Schemas {
		if err := s.ValidateSchema(schema, compare); err == nil {
			return nil
		}
	}

	return typeErrorf("was expecting %s", o.TypeString())
}

func (o *OneOf) ValidateValue(v interface{}) (interface{}, error) {
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

	if o.Const != nil && o.Const != v {
		return nil, fmt.Errorf("must be %s", o.StringValue(o.Const))
	}

	if len(o.Enum) > 0 {
		if len(o.Enum) == 1 && o.Enum[0] != v {
			return nil, fmt.Errorf("must be %s", o.StringValue(o.Enum[0]))
		}

		allowed := func() bool {
			for _, e := range o.Enum {
				if e == v {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", o.join(o.Enum, ", "))
		}
	}

	for _, s := range o.Schemas {
		vv, err := s.ValidateValue(v)
		if err == nil {
			return vv, nil
		}
	}

	return nil, fmt.Errorf("must be %s", o.TypeString())
}

func (o *OneOf) schemasString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "[]%sSchema{\n", schemaPkg(imports))
	for _, s := range o.Schemas {
		fprintf(buf, "\t%s,\n", indent(s.GoString(imports)))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (o *OneOf) join(elems []interface{}, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return o.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(o.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(o.StringValue(e))
	}
	return b.String()
}
