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
	"strings"

	"github.com/conflowio/conflow/src/internal/utils"
	"github.com/conflowio/conflow/src/util"
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
	// ParameterNames will contain the json property name -> parameter name mapping, if they are different
	// @ignore
	ParameterNames map[string]string `json:"parameterNames,omitempty"`
	// JSONPropertyNames will contain the parameter name -> json property name mapping, if they are different
	// @ignore
	JSONPropertyNames map[string]string `json:"-"`
	MinProperties     int64             `json:"minProperties,omitempty"`
	MaxProperties     *int64            `json:"maxProperties,omitempty"`
	// @name "property"
	Properties map[string]Schema `json:"properties,omitempty"`
	// Required will contain the required parameter names
	Required []string `json:"required,omitempty"`
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
	return o.Default
}

// FieldName returns the field name for the given parameter name
func (o *Object) FieldName(parameterName string) string {
	jsonPropertyName := o.JSONPropertyName(parameterName)
	fieldName := jsonPropertyName
	if name, ok := o.FieldNames[jsonPropertyName]; ok {
		fieldName = name
	}
	return fieldName
}

// JSONPropertyName returns the JSON property name for the given parameter name
func (o *Object) JSONPropertyName(parameterName string) string {
	jsonPropertyName := parameterName
	if name, ok := o.JSONPropertyNames[parameterName]; ok {
		jsonPropertyName = name
	}
	return jsonPropertyName
}

// ParameterName returns the parameter name for the given JSON property name
func (o *Object) ParameterName(jsonPropertyName string) string {
	parameterName := jsonPropertyName
	if name, ok := o.ParameterNames[jsonPropertyName]; ok {
		parameterName = name
	}
	return parameterName
}

func (o *Object) PropertyByParameterName(parameterName string) (Schema, bool) {
	s, ok := o.Properties[o.JSONPropertyName(parameterName)]
	return s, ok
}

func (o *Object) GoType(imports map[string]string) string {
	panic("GoType should not be called on object types")
}

func (o *Object) IsPropertyRequired(jsonPropertyName string) bool {
	for _, p := range o.Required {
		if p == jsonPropertyName {
			return true
		}
	}
	return false
}

func (o *Object) GoString(imports map[string]string) string {
	pkg := schemaPkg(imports)
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sObject{\n", pkg)
	if !reflect.ValueOf(o.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(o.Metadata.GoString(imports)))
	}
	if o.Const != nil {
		fprintf(buf, "\tConst: %#v,\n", o.Const)
	}
	if o.Default != nil {
		fprintf(buf, "\tDefault: %#v,\n", o.Default)
	}
	if o.DependentRequired != nil {
		fprintf(buf, "\tDependentRequired: %#v,\n", o.DependentRequired)
	}
	if len(o.Enum) > 0 {
		fprintf(buf, "\tEnum: %#v,\n", o.Enum)
	}
	if len(o.FieldNames) > 0 {
		fprintf(buf, "\tFieldNames: %#v,\n", o.FieldNames)
	}
	if len(o.JSONPropertyNames) > 0 {
		fprintf(buf, "\tJSONPropertyNames: %#v,\n", o.JSONPropertyNames)
	}
	if o.MinProperties > 0 {
		fprintf(buf, "\tMinProperties: %d,\n", o.MinProperties)
	}
	if o.MaxProperties != nil {
		fprintf(buf, "\tMaxProperties: %sPointer(int64(%d)),\n", pkg, *o.MaxProperties)
	}
	if len(o.ParameterNames) > 0 {
		fprintf(buf, "\tParameterNames: %#v,\n", o.ParameterNames)
	}
	if len(o.Properties) > 0 {
		fprintf(buf, "\tProperties: %s,\n", indent(o.propertiesString(imports)))
	}
	if len(o.Required) > 0 {
		fprintf(buf, "\tRequired: %#v,\n", o.Required)
	}
	buf.WriteRune('}')
	return buf.String()
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
	for i, p := range util.SortedMapKeys(o.Properties) {
		if i > 0 {
			sb.WriteString(", ")
		}
		fprintf(sb, "%s: %s", p, o.Properties[p].TypeString())
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

	for jsonPropertyName, fieldName := range o.FieldNames {
		if _, ok := v.Properties[jsonPropertyName]; !ok {
			return fmt.Errorf("property %q does not exist", jsonPropertyName)
		}
		if !FieldNameRegexp.MatchString(fieldName) {
			return fmt.Errorf("invalid field name %q, must match %s", fieldName, FieldNameRegexp.String())
		}
	}

	for jsonPropertyName, parameterName := range o.ParameterNames {
		if _, ok := v.Properties[jsonPropertyName]; !ok {
			return fmt.Errorf("property %q does not exist", jsonPropertyName)
		}
		if !NameRegExp.MatchString(parameterName) {
			return fmt.Errorf("invalid parameter name %q, must match %s", parameterName, NameRegExp.String())
		}
	}

	if v.Properties != nil {
		allFieldNames := make(map[string]bool, len(v.Properties))
		allParameterNames := make(map[string]bool, len(v.Properties))
		o.Properties = map[string]Schema{}
		for jsonPropertyName, s := range v.Properties {
			if s == nil {
				return fmt.Errorf("no valid schema found for %q property", jsonPropertyName)
			}
			o.Properties[jsonPropertyName] = s.Schema

			parameterName := v.ParameterNames[jsonPropertyName]
			if parameterName == "" {
				parameterName = utils.ToSnakeCase(jsonPropertyName)
			}

			if allParameterNames[parameterName] {
				return fmt.Errorf("multiple properties are using the %q parameter name", parameterName)
			}
			allParameterNames[parameterName] = true

			if parameterName != jsonPropertyName {
				if o.JSONPropertyNames == nil {
					o.JSONPropertyNames = map[string]string{}
				}
				o.JSONPropertyNames[parameterName] = jsonPropertyName
			}

			fieldName := v.FieldNames[jsonPropertyName]
			if fieldName == "" {
				if !FieldNameRegexp.MatchString(jsonPropertyName) {
					return fmt.Errorf("property name %q can not be used as a field name", jsonPropertyName)
				}
				fieldName = jsonPropertyName
			}
			if allFieldNames[fieldName] {
				return fmt.Errorf("multiple properties are using the %q field name", fieldName)
			}
			allFieldNames[fieldName] = true
		}
	}

	return nil
}

func (o *Object) Validate(ctx *Context) error {
	return ValidateAll(
		ctx,
		validateCommonFields(o, o.Const, o.Default, o.Enum),
		ValidateMap("properties", o.Properties),
	)
}

func (o *Object) ValidateSchema(s Schema, compare bool) error {
	if s.Type() == TypeNull {
		return nil
	}

	o2, ok := s.(*Object)
	if !ok {
		return typeError("must be object")
	}

	for n2, p2 := range o2.Properties {
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
		if _, ok := o2.Properties[required]; !ok {
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
		if p, ok := o.Properties[k]; ok {
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

func (o *Object) propertiesString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "map[string]%sSchema{\n", schemaPkg(imports))
	for _, k := range util.SortedMapKeys(o.Properties) {
		if p := o.Properties[k]; p != nil {
			fprintf(buf, "\t%#v: %s,\n", k, indent(p.GoString(imports)))
		}
	}
	buf.WriteRune('}')
	return buf.String()
}
