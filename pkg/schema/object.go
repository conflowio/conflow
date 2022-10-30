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
	"text/template"

	"github.com/conflowio/conflow/pkg/internal/utils"
	"github.com/conflowio/conflow/pkg/util"
	"github.com/conflowio/conflow/pkg/util/validation"
)

type structTemplateParams struct {
	Name    string
	Object  *Object
	Imports map[string]string
}

const structTemplate = `{{ $root := . -}}
struct {
{{ range $jsonPropertyName, $p := .Object.Properties -}}
	{{ with $root.Object.ParameterName $jsonPropertyName }}{{ if not (eq . $jsonPropertyName) }}// @name "{{ . }}"
{{ end }}{{ end -}}
	{{ with description ($root.Object.FieldName $jsonPropertyName) $p }}{{ . }}{{ end -}}
	{{ $root.Object.FieldName $jsonPropertyName}} {{ $p.GoType $root.Imports }}{{ with structTags $jsonPropertyName}} {{ . }}{{ end }}
{{ end -}}
}`

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type Object struct {
	Metadata

	Const             map[string]interface{}   `json:"const,omitempty"`
	Default           map[string]interface{}   `json:"default,omitempty"`
	DependentRequired map[string][]string      `json:"dependentRequired,omitempty"`
	Enum              []map[string]interface{} `json:"enum,omitempty"`
	// FieldNames will contain the json property name -> field name mapping, if they are different
	FieldNames map[string]string `json:"x-conflow-fields,omitempty"`
	// ParameterNames will contain the json property name -> parameter name mapping, if they are different
	ParameterNames map[string]string `json:"x-conflow-parameters,omitempty"`
	MinProperties  int64             `json:"minProperties,omitempty"`
	MaxProperties  *int64            `json:"maxProperties,omitempty"`
	Nullable       bool              `json:"nullable,omitempty"`
	// @name "property"
	Properties map[string]Schema `json:"properties,omitempty"`
	// Required will contain the required parameter names
	Required []string `json:"required,omitempty"`

	// @ignore
	jsonPropertyNames map[string]string
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

func (o *Object) GetNullable() bool {
	return o.Nullable
}

// JSONPropertyName returns the JSON property name for the given parameter name
func (o *Object) JSONPropertyName(parameterName string) string {
	if o.jsonPropertyNames == nil {
		o.jsonPropertyNames = make(map[string]string, len(o.ParameterNames))
		for prop, param := range o.ParameterNames {
			o.jsonPropertyNames[param] = prop
		}
	}

	jsonPropertyName := parameterName
	if name, ok := o.jsonPropertyNames[parameterName]; ok {
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

func (o *Object) SetNullable(nullable bool) {
	o.Nullable = nullable
}

func (o *Object) GoType(imports map[string]string) string {
	id := o.GetID()
	if id != "" {
		parts := strings.Split(id, ".")
		name := parts[len(parts)-1]
		pkg := strings.Join(parts[0:len(parts)-1], ".")
		sel := utils.EnsureUniqueGoPackageSelector(imports, pkg)
		if o.Nullable {
			return fmt.Sprintf("*%s%s", sel, name)
		}
		return fmt.Sprintf("%s%s", sel, name)
	}

	res, err := o.GenerateStruct(imports)
	if err != nil {
		panic(fmt.Errorf("failed to generate struct for object: %w", err))
	}
	return string(res)
}

func (o *Object) GenerateStruct(imports map[string]string) ([]byte, error) {
	return util.GenerateTemplate(
		structTemplate,
		structTemplateParams{
			Object:  o,
			Imports: imports,
		},
		template.FuncMap{
			"description": func(name string, s Schema) string {
				description := s.GetDescription()
				if description == "" {
					return ""
				}
				if strings.HasPrefix(description, "It ") {
					description = strings.Replace(description, "It ", "", 1)
				}
				return fmt.Sprintf("// %s %s\n", name, strings.ReplaceAll(description, "\n", "\n//"))
			},
			"structTags": func(propertyName string) string {
				fieldName, ok := o.FieldNames[propertyName]
				if !ok {
					fieldName = propertyName
				}
				if fieldName == propertyName {
					return ""
				}
				omitEmptyStr := ",omitempty"
				if o.IsPropertyRequired(propertyName) {
					omitEmptyStr = ""
				}
				return fmt.Sprintf("`json:\"%s%s\"`", propertyName, omitEmptyStr)
			},
		},
	)
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
		func(ctx *Context) error {
			if o.MinProperties < 0 {
				return validation.NewFieldError("minProperties", errors.New("must be greater than or equal to 0"))
			}
			if o.MaxProperties != nil && o.MinProperties > *o.MaxProperties {
				return errors.New("minProperties and maxProperties constraints are impossible to fulfil")
			}
			if int64(len(o.Properties)) < o.MinProperties {
				return validation.NewFieldError("minProperties", errors.New("can not be greater than the number of properties defined"))
			}

			for i, r := range o.Required {
				if _, ok := o.Properties[r]; !ok {
					return validation.NewFieldErrorf(fmt.Sprintf("required[%d]", i), "property %q does not exist", r)
				}
			}

			return nil
		},
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
	var ve validation.Error

	m, ok := value.(map[string]interface{})
	if !ok {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return nil, fmt.Errorf("can not be empty")
			} else {
				v = v.Elem()
				value = v.Interface()
			}
		}
		t := v.Type()
		if t.Kind() != reflect.Struct {
			return nil, fmt.Errorf("must be object")
		}

		m = map[string]interface{}{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if !f.IsExported() {
				continue
			}

			name := f.Name
			var omitEmpty bool
			if jsonTag, ok := f.Tag.Lookup("json"); ok {
				parts := strings.Split(jsonTag, ",")
				if parts[0] == "-" {
					continue
				}
				if parts[0] != "" {
					name = parts[0]
				}
				if len(parts) > 1 && parts[1] == "omitempty" {
					omitEmpty = true
				}
			}

			vf := v.Field(i)

			if p, ok := o.Properties[name]; ok {
				if vf.Kind() == reflect.Interface && vf.IsZero() {
					continue
				}
				validatedValue, err := p.ValidateValue(v.Field(i).Interface())
				if err != nil {
					ve.AddFieldError(name, err)
				}
				if !omitEmpty || !v.Field(i).IsZero() {
					m[name] = validatedValue
				}
			} else {
				ve.AddFieldError(name, errors.New("property does not exist"))
			}
		}
	} else {
		for name, v := range m {
			if p, ok := o.Properties[name]; ok {
				fv, err := p.ValidateValue(v)
				if err != nil {
					ve.AddFieldError(name, err)
				}
				m[name] = fv
			} else {
				ve.AddFieldError(name, errors.New("property does not exist"))
			}
		}
	}

	if o.Const != nil && o.CompareValues(o.Const, m) != 0 {
		return nil, fmt.Errorf("must be %s", o.StringValue(o.Const))
	}

	if len(o.Enum) == 1 && o.CompareValues(o.Enum[0], m) != 0 {
		return nil, fmt.Errorf("must be %s", o.StringValue(o.Enum[0]))
	}

	if len(o.Enum) > 0 {
		allowed := func() bool {
			for _, e := range o.Enum {
				if o.CompareValues(e, m) == 0 {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", o.join(o.Enum, ", "))
		}
	}

	if int64(len(m)) < o.MinProperties {
		switch o.MinProperties {
		case 1:
			ve.AddError(errors.New("the object can not be empty"))
		default:
			ve.AddErrorf("the object must have at least %d properties defined", o.MinProperties)
		}
	}

	if o.MaxProperties != nil && int64(len(m)) > *o.MaxProperties {
		switch *o.MaxProperties {
		case 0:
			ve.AddError(errors.New("the object must be empty"))
		case 1:
			ve.AddError(errors.New("the object can only have a single property defined"))
		default:
			ve.AddErrorf("the object can not have more than %d properties defined", *o.MaxProperties)
		}
	}

	if len(o.Required) > 0 || len(o.DependentRequired) > 0 {
		missingFields := map[string]bool{}

		for _, f := range o.Required {
			if _, ok := m[f]; !ok {
				missingFields[f] = true
			}
		}

		for p, required := range o.DependentRequired {
			if _, ok := m[p]; !ok {
				continue
			}
			for _, f := range required {
				if _, ok := m[f]; !ok {
					missingFields[f] = true
				}
			}
		}

		for f := range missingFields {
			ve.AddFieldError(f, errors.New("required"))
		}
	}

	if err := ve.ErrOrNil(); err != nil {
		return nil, err
	}

	if o.Nullable {
		return &value, nil
	}
	return value, nil
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
