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
	"sort"
	"strings"
)

type Object struct {
	Metadata

	Const            *map[string]interface{}  `json:"const,omitempty"`
	Default          *map[string]interface{}  `json:"default,omitempty"`
	Enum             []map[string]interface{} `json:"enum,omitempty"`
	Properties       map[string]Schema        `json:"properties,omitempty"`
	Required         []string                 `json:"required,omitempty"`
	StructProperties map[string]string        `json:"structProperties,omitempty"`
}

func (o *Object) AssignValue(_ map[string]string, _, _ string) string {
	panic("AssignValue should not be called on an object type")
}

func (o *Object) CompareValues(v1, v2 interface{}) int {
	o1, ok := v1.(map[string]interface{})
	if !ok {
		return -1
	}

	o2, ok := v2.(map[string]interface{})
	if !ok {
		return 1
	}

	switch {
	case len(o1) == len(o2):
		keys1 := getSortedMapKeys(o1)
		keys2 := getSortedMapKeys(o2)

		for i := 0; i < len(o1); i++ {
			k1 := keys1[i]
			k2 := keys2[i]
			switch {
			case k1 == k2:
				if c := o.Properties[k1].CompareValues(o1[k1], o2[k2]); c != 0 {
					return c
				}
			case k1 < k2:
				return -1
			default:
				return 1
			}
		}

		return 0
	case len(o1) < len(o2):
		return -1
	default:
		return 1
	}
}

func (o *Object) Copy() Schema {
	j, err := json.Marshal(o)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Object{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (o *Object) DefaultValue() interface{} {
	if o.Default == nil {
		return nil
	}
	return *o.Default
}

func (o *Object) GetProperties() map[string]Schema {
	return o.Properties
}

func (o *Object) GetRequired() []string {
	return o.Required
}

func (o *Object) GetStructProperties() map[string]string {
	return o.StructProperties
}

func (o *Object) GoType(imports map[string]string) string {
	panic("GoType should not be called on object types")
}

func (o *Object) IsPropertyRequired(name string) bool {
	for _, p := range o.Required {
		if p == name {
			return true
		}
	}
	return false
}

func (o *Object) MarshalJSON() ([]byte, error) {
	type Alias Object
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(o.Type()),
		Alias: (*Alias)(o),
	})
}

func (o *Object) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Object{\n")
	if !reflect.ValueOf(o.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(o.Metadata.GoString()))
	}
	if o.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: %#v,\n", o.Const)
	}
	if o.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: %#v,\n", o.Default)
	}
	if len(o.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", o.Enum)
	}
	if len(o.Properties) > 0 {
		_, _ = fmt.Fprintf(buf, "\tProperties: %s,\n", indent(o.propertiesString()))
	}
	if len(o.Required) > 0 {
		_, _ = fmt.Fprintf(buf, "\tRequired: %#v,\n", o.Required)
	}
	if len(o.StructProperties) > 0 {
		_, _ = fmt.Fprintf(buf, "\tStructProperties: %#v,\n", o.StructProperties)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (o *Object) StringValue(value interface{}) string {
	v, ok := value.(map[string]interface{})
	if !ok {
		return ""
	}

	if len(v) == 0 {
		return "{}"
	}

	keys := getSortedMapKeys(v)

	var b strings.Builder
	b.WriteRune('{')
	b.WriteString(keys[0])
	b.WriteString(": ")
	b.WriteString(o.Properties[keys[0]].StringValue(v[keys[0]]))
	for _, k := range keys[1:] {
		b.WriteString(", ")
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(o.Properties[k].StringValue(v[k]))
	}
	b.WriteRune('}')
	return b.String()
}

func (o *Object) Type() Type {
	return TypeObject
}

func (o *Object) TypeString() string {
	sb := &strings.Builder{}
	sb.WriteString("object(")
	for i, p := range o.propertyNames() {
		if i > 0 {
			sb.WriteString(", ")
		}
		_, _ = fmt.Fprintf(sb, "%s: %s", p, o.Properties[p].TypeString())
	}
	sb.WriteRune(')')
	return sb.String()
}

func (o *Object) UnmarshalJSON(j []byte) error {
	type Alias Object
	v := struct {
		*Alias
		Properties map[string]*SchemaUnmarshaler `json:"properties,omitempty"`
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	if v.Properties != nil {
		o.Properties = map[string]Schema{}
		for p, s := range v.Properties {
			if s != nil {
				o.Properties[p] = s.Schema
			}
		}
	}

	return nil
}

func (o *Object) ValidateSchema(s Schema, compare bool) error {
	if s.Type() == TypeNull {
		return nil
	}

	o2, ok := s.(ObjectKind)
	if !ok {
		return typeError("must be object")
	}

	for n2, p2 := range o2.GetProperties() {
		p, ok := o.Properties[n2]
		if !ok {
			return typeErrorf("was expecting %s", o.TypeString())
		}

		if err := p.ValidateSchema(p2, compare); err != nil {
			if _, ok := err.(typeError); ok {
				return typeErrorf("was expecting %s", o.TypeString())
			}
			return err
		}
	}

	for _, required := range o.Required {
		if _, ok := o2.GetProperties()[required]; !ok {
			return typeErrorf("was expecting %s", o.TypeString())
		}
	}

	return nil
}

func (o *Object) ValidateValue(value interface{}) error {
	v, ok := value.(map[string]interface{})
	if !ok {
		return errors.New("must be object")
	}

	if o.Const != nil && o.CompareValues(*o.Const, v) != 0 {
		return fmt.Errorf("must be %s", o.StringValue(*o.Const))
	}

	if len(o.Enum) == 1 && o.CompareValues(o.Enum[0], v) != 0 {
		return fmt.Errorf("must be %s", o.StringValue(o.Enum[0]))
	}

	if len(o.Enum) > 0 {
		allowed := func() bool {
			for _, e := range o.Enum {
				if o.CompareValues(e, v) == 0 {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return fmt.Errorf("must be one of %s", o.join(o.Enum, ", "))
		}
	}

	for _, f := range o.Required {
		if _, ok := v[f]; !ok {
			return NewFieldError(f, errors.New("required"))
		}
	}

	for _, k := range getSortedMapKeys(v) {
		if p, ok := o.Properties[k]; ok {
			if err := p.ValidateValue(v[k]); err != nil {
				return NewFieldError(k, err)
			}
		} else {
			return NewFieldError(k, errors.New("property does not exist"))
		}
	}

	return nil
}

func (o *Object) join(elems []map[string]interface{}, sep string) string {
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

func (o *Object) propertiesString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("map[string]schema.Schema{\n")
	for _, k := range o.propertyNames() {
		if p := o.Properties[k]; p != nil {
			_, _ = fmt.Fprintf(buf, "\t%#v: %s,\n", k, indent(p.GoString()))
		}
	}
	buf.WriteRune('}')
	return buf.String()
}

func (o *Object) propertyNames() []string {
	if len(o.Properties) == 0 {
		return nil
	}

	keys := make([]string, 0, len(o.Properties))
	for k := range o.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func ObjectPtr(v map[string]interface{}) *map[string]interface{} {
	return &v
}
