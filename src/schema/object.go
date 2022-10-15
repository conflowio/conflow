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

	"github.com/conflowio/conflow/src/internal/utils"
)

type Object struct {
	Metadata

	Const             map[string]interface{}   `json:"const,omitempty"`
	Default           map[string]interface{}   `json:"default,omitempty"`
	DependentRequired map[string][]string      `json:"dependentRequired,omitempty"`
	Enum              []map[string]interface{} `json:"enum,omitempty"`
	// FieldNames will contain the json property name -> field name mapping, if they are different
	// @ignore
	FieldNames map[string]string `json:"fieldNames,omitempty"`
	// JSONPropertyNames will contain the parameter name -> json property name mapping, if they are different
	// @ignore
	JSONPropertyNames map[string]string `json:"-"`
	MinProperties     int64             `json:"minProperties,omitempty"`
	MaxProperties     *int64            `json:"maxProperties,omitempty"`
	// @name "property"
	Parameters map[string]Schema `json:"-"`
	// Required will contain the required parameter names
	// @ignore
	Required []string `json:"-"`
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
				if c := o.Parameters[k1].CompareValues(o1[k1], o2[k2]); c != 0 {
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
	return o.Default
}

// GetFieldName returns the field name for the given parameter name
func (o *Object) GetFieldName(parameterName string) string {
	jsonPropertyName := o.GetJSONPropertyName(parameterName)
	fieldName := jsonPropertyName
	if name, ok := o.FieldNames[jsonPropertyName]; ok {
		fieldName = name
	}
	return fieldName
}

// GetJSONPropertyName returns the JSON property name for the given parameter name
func (o *Object) GetJSONPropertyName(parameterName string) string {
	jsonPropertyName := parameterName
	if name, ok := o.JSONPropertyNames[parameterName]; ok {
		jsonPropertyName = name
	}
	return jsonPropertyName
}

func (o *Object) GetParameters() map[string]Schema {
	return o.Parameters
}

func (o *Object) GetRequired() []string {
	return o.Required
}

func (o *Object) GoType(imports map[string]string) string {
	panic("GoType should not be called on object types")
}

func (o *Object) IsParameterRequired(name string) bool {
	for _, p := range o.Required {
		if p == name {
			return true
		}
	}
	return false
}

func (o *Object) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Object{\n")
	if !reflect.ValueOf(o.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(o.Metadata.GoString(imports)))
	}
	if o.Const != nil {
		_, _ = fmt.Fprintf(buf, "\tConst: %#v,\n", o.Const)
	}
	if o.Default != nil {
		_, _ = fmt.Fprintf(buf, "\tDefault: %#v,\n", o.Default)
	}
	if o.DependentRequired != nil {
		_, _ = fmt.Fprintf(buf, "\tDependentRequired: %#v,\n", o.DependentRequired)
	}
	if len(o.Enum) > 0 {
		_, _ = fmt.Fprintf(buf, "\tEnum: %#v,\n", o.Enum)
	}
	if len(o.FieldNames) > 0 {
		_, _ = fmt.Fprintf(buf, "\tFieldNames: %#v,\n", o.FieldNames)
	}
	if len(o.JSONPropertyNames) > 0 {
		_, _ = fmt.Fprintf(buf, "\tJSONPropertyNames: %#v,\n", o.JSONPropertyNames)
	}
	if o.MinProperties > 0 {
		_, _ = fmt.Fprintf(buf, "\tMinProperties: %d,\n", o.MinProperties)
	}
	if o.MaxProperties != nil {
		_, _ = fmt.Fprintf(buf, "\tMaxProperties: schema.Pointer(int64(%d)),\n", *o.MaxProperties)
	}
	if len(o.Parameters) > 0 {
		_, _ = fmt.Fprintf(buf, "\tParameters: %s,\n", indent(o.parametersString(imports)))
	}
	if len(o.Required) > 0 {
		_, _ = fmt.Fprintf(buf, "\tRequired: %#v,\n", o.Required)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (o *Object) MarshalJSON() ([]byte, error) {
	parameterNames := map[string]string{}
	properties := map[string]Schema{}
	required := make([]string, 0, len(o.Required))
	for parameterName, schema := range o.Parameters {
		jsonPropertyName := o.GetJSONPropertyName(parameterName)
		if jsonPropertyName != parameterName {
			parameterNames[jsonPropertyName] = parameterName
		}
		properties[jsonPropertyName] = schema
		if o.IsParameterRequired(parameterName) {
			required = append(required, jsonPropertyName)
		}
	}

	type Alias Object
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
		ParameterNames map[string]string `json:"parameterNames,omitempty"`
		Properties     map[string]Schema `json:"properties,omitempty"`
		Required       []string          `json:"required,omitempty"`
	}{
		Type:           string(o.Type()),
		Alias:          (*Alias)(o),
		ParameterNames: parameterNames,
		Properties:     properties,
		Required:       required,
	})
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
	b.WriteString(o.Parameters[keys[0]].StringValue(v[keys[0]]))
	for _, k := range keys[1:] {
		b.WriteString(", ")
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(o.Parameters[k].StringValue(v[k]))
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
	for i, p := range o.parameterNames() {
		if i > 0 {
			sb.WriteString(", ")
		}
		_, _ = fmt.Fprintf(sb, "%s: %s", p, o.Parameters[p].TypeString())
	}
	sb.WriteRune(')')
	return sb.String()
}

func (o *Object) UnmarshalJSON(j []byte) error {
	type Alias Object
	v := struct {
		*Alias
		Properties     map[string]*SchemaUnmarshaler `json:"properties,omitempty"`
		ParameterNames map[string]string             `json:"parameterNames,omitempty"`
		Required       []string                      `json:"required,omitempty"`
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	getParameterName := func(jsonPropertyName string) string {
		if name, ok := v.ParameterNames[jsonPropertyName]; ok && NameRegExp.MatchString(name) {
			return name
		} else {
			return utils.ToSnakeCase(jsonPropertyName)
		}
	}

	if v.Properties != nil {
		o.Parameters = map[string]Schema{}
		for jsonPropertyName, su := range v.Properties {
			if su != nil {
				o.Parameters[getParameterName(jsonPropertyName)] = su.Schema
			}
			parameterName := getParameterName(jsonPropertyName)
			if parameterName != jsonPropertyName {
				if o.JSONPropertyNames == nil {
					o.JSONPropertyNames = map[string]string{}
				}
				o.JSONPropertyNames[parameterName] = jsonPropertyName
			}
		}
	}

	if len(v.Required) > 0 {
		o.Required = make([]string, 0, len(v.Required))
		for _, jsonPropertyName := range v.Required {
			o.Required = append(o.Required, getParameterName(jsonPropertyName))
		}
	}

	for jsonPropertyName, fieldName := range o.FieldNames {
		if !fieldNameRegexp.MatchString(fieldName) {
			delete(o.FieldNames, jsonPropertyName)
		}
	}
	if len(o.FieldNames) == 0 {
		o.FieldNames = nil
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

	for n2, p2 := range o2.GetParameters() {
		p, ok := o.Parameters[n2]
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
		if _, ok := o2.GetParameters()[required]; !ok {
			return typeErrorf("was expecting %s", o.TypeString())
		}
	}

	return nil
}

func (o *Object) ValidateValue(value interface{}) (interface{}, error) {
	v, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New("must be object")
	}

	if o.Const != nil && o.CompareValues(o.Const, v) != 0 {
		return nil, fmt.Errorf("must be %s", o.StringValue(o.Const))
	}

	if len(o.Enum) == 1 && o.CompareValues(o.Enum[0], v) != 0 {
		return nil, fmt.Errorf("must be %s", o.StringValue(o.Enum[0]))
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
			return nil, fmt.Errorf("must be one of %s", o.join(o.Enum, ", "))
		}
	}

	if int64(len(v)) < o.MinProperties {
		switch o.MinProperties {
		case 1:
			return nil, errors.New("the object can not be empty")
		default:
			return nil, fmt.Errorf("the object must have at least %d properties defined", o.MinProperties)
		}
	}

	if o.MaxProperties != nil && int64(len(v)) > *o.MaxProperties {
		switch *o.MaxProperties {
		case 0:
			return nil, errors.New("the object must be empty")
		case 1:
			return nil, errors.New("the object can only have a single property defined")
		default:
			return nil, fmt.Errorf("the object can not have more than %d properties defined", *o.MaxProperties)
		}
	}

	ve := ValidationError{}

	if len(o.Required) > 0 || len(o.DependentRequired) > 0 {
		missingFields := map[string]bool{}

		for _, f := range o.Required {
			if _, ok := v[f]; !ok {
				missingFields[f] = true
			}
		}

		for p, required := range o.DependentRequired {
			if _, ok := v[p]; !ok {
				continue
			}
			for _, f := range required {
				if _, ok := v[f]; !ok {
					missingFields[f] = true
				}
			}
		}

		for f := range missingFields {
			ve.AddError(f, errors.New("required"))
		}
	}

	for _, k := range getSortedMapKeys(v) {
		if p, ok := o.Parameters[k]; ok {
			nv, err := p.ValidateValue(v[k])
			if err != nil {
				ve.AddError(k, err)
			} else {
				v[k] = nv
			}
		} else {
			ve.AddError(k, errors.New("property does not exist"))
		}
	}

	if err := ve.ErrOrNil(); err != nil {
		return nil, err
	}

	return v, nil
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

func (o *Object) parametersString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("map[string]schema.Schema{\n")
	for _, k := range o.parameterNames() {
		if p := o.Parameters[k]; p != nil {
			_, _ = fmt.Fprintf(buf, "\t%#v: %s,\n", k, indent(p.GoString(imports)))
		}
	}
	buf.WriteRune('}')
	return buf.String()
}

func (o *Object) parameterNames() []string {
	if len(o.Parameters) == 0 {
		return nil
	}

	keys := make([]string, 0, len(o.Parameters))
	for k := range o.Parameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
