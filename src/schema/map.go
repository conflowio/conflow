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
	"regexp"
	"strings"
)

type Map struct {
	Metadata

	// @required
	AdditionalProperties Schema                   `json:"additionalProperties,omitempty"`
	Const                map[string]interface{}   `json:"const,omitempty"`
	Default              map[string]interface{}   `json:"default,omitempty"`
	Enum                 []map[string]interface{} `json:"enum,omitempty"`
	MinProperties        int64                    `json:"minProperties,omitempty"`
	MaxProperties        *int64                   `json:"maxProperties,omitempty"`
}

func (m *Map) AssignValue(imports map[string]string, valueName, resultName string) string {
	if m.AdditionalProperties.Type() == TypeAny {
		return fmt.Sprintf("%s = %s.(map[string]interface{})", resultName, valueName)
	}

	valueNameVar := regexp.MustCompile(`[\[\]]`).ReplaceAllString(valueName, "")

	return fmt.Sprintf(`%s = make(map[string]%s, len(%s.(map[string]interface{})))
for %sk, %sv := range %s.(map[string]interface{}) {
	%s
}`,
		resultName,
		m.AdditionalProperties.GoType(imports),
		valueName,
		valueNameVar,
		valueNameVar,
		valueName,
		indent(m.AdditionalProperties.AssignValue(imports, valueNameVar+"v", fmt.Sprintf("%s[%sk]", resultName, valueNameVar))),
	)
}

func (m *Map) CompareValues(v1, v2 interface{}) int {
	var m1 map[string]interface{}
	if v1 != nil {
		var ok bool
		if m1, ok = v1.(map[string]interface{}); !ok {
			return -1
		}
	}

	var m2 map[string]interface{}
	if v2 != nil {
		var ok bool
		if m2, ok = v2.(map[string]interface{}); !ok {
			return 1
		}
	}

	p := m.GetAdditionalProperties()

	switch {
	case len(m1) == len(m2):
		keys1 := getSortedMapKeys(m1)
		keys2 := getSortedMapKeys(m2)

		for i := 0; i < len(m1); i++ {
			k1 := keys1[i]
			k2 := keys2[i]
			switch {
			case k1 == k2:
				if c := p.CompareValues(m1[k1], m2[k2]); c != 0 {
					return c
				}
			case k1 < k2:
				return -1
			default:
				return 1
			}
		}

		return 0
	case len(m1) < len(m2):
		return -1
	default:
		return 1
	}
}

func (m *Map) Copy() Schema {
	j, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Map{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (m *Map) DefaultValue() interface{} {
	if m.Default == nil {
		return nil
	}

	return m.Default
}

func (m *Map) GetAdditionalProperties() Schema {
	if m.AdditionalProperties == nil {
		return AnyValue()
	}

	return m.AdditionalProperties
}

func (m *Map) GoType(imports map[string]string) string {
	return fmt.Sprintf("map[string]%s", m.GetAdditionalProperties().GoType(imports))
}

func (m *Map) MarshalJSON() ([]byte, error) {
	type Alias Map
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(TypeObject),
		Alias: (*Alias)(m),
	})
}

func (m *Map) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Map{\n")
	if !reflect.ValueOf(m.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(m.Metadata.GoString(imports)))
	}
	if m.AdditionalProperties != nil {
		_, _ = fmt.Fprintf(buf, "\tAdditionalProperties: %s,\n", indent(m.AdditionalProperties.GoString(imports)))
	}
	if m.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: %#v,\n", m.Const)
	}
	if m.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: %#v,\n", m.Default)
	}
	if len(m.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", m.Enum)
	}
	if m.MinProperties > 0 {
		_, _ = fmt.Fprintf(buf, "\tMinProperties: %d,\n", m.MinProperties)
	}
	if m.MaxProperties != nil {
		_, _ = fmt.Fprintf(buf, "\tMaxProperties: schema.Pointer(int64(%d)),\n", *m.MaxProperties)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (m *Map) StringValue(value interface{}) string {
	v, ok := value.(map[string]interface{})
	if !ok {
		return ""
	}

	if len(v) == 0 {
		return "map{}"
	}

	keys := getSortedMapKeys(v)
	p := m.GetAdditionalProperties()

	var b strings.Builder
	b.WriteString("map{")
	b.WriteString(keys[0])
	b.WriteString(": ")
	b.WriteString(p.StringValue(v[keys[0]]))
	for _, k := range keys[1:] {
		b.WriteString(", ")
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(p.StringValue(v[k]))
	}
	b.WriteRune('}')
	return b.String()
}

func (m *Map) Type() Type {
	return TypeMap
}

func (m *Map) TypeString() string {
	if m.AdditionalProperties == nil || m.AdditionalProperties.Type() == TypeFalse {
		return "map"
	}

	return fmt.Sprintf("map(%s)", m.AdditionalProperties.TypeString())
}

func (m *Map) UnmarshalJSON(j []byte) error {
	type Alias Map
	v := struct {
		*Alias
		AdditionalProperties *SchemaUnmarshaler `json:"additionalProperties,omitempty"`
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	if v.AdditionalProperties != nil {
		m.AdditionalProperties = v.AdditionalProperties.Schema
	}

	return nil
}

func (m *Map) Validate(ctx *Context) error {
	return ValidateAll(
		ctx,
		validateCommonFields(m, m.Const, m.Default, m.Enum),
		Validate("additionalProperties", m.AdditionalProperties),
	)
}

func (m *Map) ValidateSchema(s Schema, compare bool) error {
	if s.Type() == TypeNull {
		return nil
	}

	o2, ok := s.(*Map)
	if !ok {
		return typeError("must be map")
	}

	switch {
	case isTypedMap(m) && !isTypedMap(o2):
		return typeErrorf("was expecting %s", m.TypeString())
	case isTypedMap(m) && isTypedMap(o2):
		if err := m.GetAdditionalProperties().ValidateSchema(o2.GetAdditionalProperties(), compare); err != nil {
			if _, ok := err.(typeError); ok {
				return typeErrorf("was expecting %s", m.TypeString())
			}
			return err
		}
	}

	return nil
}

func (m *Map) ValidateValue(value interface{}) (interface{}, error) {
	var v map[string]interface{}
	if value != nil {
		var ok bool
		if v, ok = value.(map[string]interface{}); !ok {
			return nil, errors.New("must be map")
		}
	}

	if m.Const != nil && m.CompareValues(m.Const, v) != 0 {
		return nil, fmt.Errorf("must be %s", m.StringValue(m.Const))
	}

	if len(m.Enum) == 1 && m.CompareValues(m.Enum[0], v) != 0 {
		return nil, fmt.Errorf("must be %s", m.StringValue(m.Enum[0]))
	}

	if len(m.Enum) > 0 {
		allowed := func() bool {
			for _, e := range m.Enum {
				if m.CompareValues(e, v) == 0 {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", m.join(m.Enum, ", "))
		}
	}

	if int64(len(v)) < m.MinProperties {
		switch m.MinProperties {
		case 1:
			return nil, errors.New("the map can not be empty")
		default:
			return nil, fmt.Errorf("the map must contain at least %d elements", m.MinProperties)
		}
	}

	if m.MaxProperties != nil && int64(len(v)) > *m.MaxProperties {
		switch *m.MaxProperties {
		case 0:
			return nil, errors.New("the map must be empty")
		case 1:
			return nil, errors.New("the map can only have a single element")
		default:
			return nil, fmt.Errorf("the map can not have more than %d elements", *m.MaxProperties)
		}
	}

	p := m.GetAdditionalProperties()

	for _, k := range getSortedMapKeys(v) {
		if p.Type() == TypeFalse {
			return nil, NewFieldError(k, errors.New("no map values are allowed"))
		} else {
			nv, err := p.ValidateValue(v[k])
			if err != nil {
				return nil, NewFieldError(k, err)
			} else {
				v[k] = nv
			}
		}
	}

	return v, nil
}

func (m *Map) join(elems []map[string]interface{}, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return m.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(m.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(m.StringValue(e))
	}
	return b.String()
}

func isTypedMap(o *Map) bool {
	return o.GetAdditionalProperties().Type() != TypeAny
}
