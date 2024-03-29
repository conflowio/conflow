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
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/conflowio/conflow/pkg/util/validation"
)

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type Array struct {
	Metadata

	Const   []interface{}   `json:"const,omitempty"`
	Default []interface{}   `json:"default,omitempty"`
	Enum    [][]interface{} `json:"enum,omitempty"`
	// @required
	Items       Schema `json:"items,omitempty"`
	MinItems    int64  `json:"minItems,omitempty"`
	MaxItems    *int64 `json:"maxItems,omitempty"`
	UniqueItems bool   `json:"uniqueItems,omitempty"`
}

func (a *Array) AssignValue(imports map[string]string, valueName, resultName string) string {
	if a.Items.Type() == TypeAny {
		return fmt.Sprintf("%s = %s.([]interface{})", resultName, valueName)
	}

	valueNameVar := regexp.MustCompile(`[\[\]]`).ReplaceAllString(valueName, "")

	return fmt.Sprintf(`%s = make(%s, len(%s.([]interface{})))
for %sk, %sv := range %s.([]interface{}) {
	%s
}`,
		resultName,
		a.GoType(imports),
		valueName,
		valueNameVar,
		valueNameVar,
		valueName,
		indent(a.Items.AssignValue(imports, valueNameVar+"v", fmt.Sprintf("%s[%sk]", resultName, valueNameVar))),
	)
}

func (a *Array) CompareValues(v1, v2 interface{}) int {
	var a1 []interface{}
	if v1 != nil {
		var ok bool
		if a1, ok = v1.([]interface{}); !ok {
			return -1
		}
	}

	var a2 []interface{}
	if v2 != nil {
		var ok bool
		if a2, ok = v2.([]interface{}); !ok {
			return 1
		}
	}

	switch {
	case len(a1) == len(a2):
		for i := 0; i < len(a1); i++ {
			if c := a.Items.CompareValues(a1[i], a2[i]); c != 0 {
				return c
			}
		}
		return 0
	case len(a1) < len(a2):
		return -1
	default:
		return 1
	}
}

func (a *Array) Copy() Schema {
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

func (a *Array) DefaultValue() interface{} {
	if a.Default == nil {
		return nil
	}

	return a.Default
}

func (a *Array) GetItems() Schema {
	return a.Items
}

func (a *Array) GoType(imports map[string]string) string {
	return fmt.Sprintf("[]%s", a.Items.GoType(imports))
}

func (a *Array) GoString(imports map[string]string) string {
	pkg := schemaPkg(imports)
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sArray{\n", pkg)
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
	if a.Items != nil {
		fprintf(buf, "\tItems: %s,\n", indent(a.Items.GoString(imports)))
	}
	if a.MinItems != 0 {
		fprintf(buf, "\tMinItems: %d,\n", a.MinItems)
	}
	if a.MaxItems != nil {
		fprintf(buf, "\tMaxItems: %sPointer(int64(%d)),\n", pkg, *a.MaxItems)
	}
	if a.UniqueItems {
		_, _ = fmt.Fprint(buf, "\tUniqueItems: true,\n")
	}
	buf.WriteRune('}')
	return buf.String()
}

func (a *Array) MarshalJSON() ([]byte, error) {
	type Alias Array
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(a.Type()),
		Alias: (*Alias)(a),
	})
}

func (a *Array) StringValue(value interface{}) string {
	v, ok := value.([]interface{})
	if !ok {
		return ""
	}

	if len(v) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteRune('[')
	sb.WriteString(a.Items.StringValue(v[0]))
	for _, e := range v[1:] {
		sb.WriteString(", ")
		sb.WriteString(a.Items.StringValue(e))
	}
	sb.WriteRune(']')
	return sb.String()
}

func (a *Array) Type() Type {
	return TypeArray
}

func (a *Array) TypeString() string {
	if a.Items.Type() == TypeAny {
		return "array"
	}

	return fmt.Sprintf("array(%s)", a.Items.TypeString())
}

func (a *Array) UnmarshalJSON(j []byte) error {
	type Alias Array
	v := struct {
		*Alias
		Items *SchemaUnmarshaler `json:"items,omitempty"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	a.Items = v.Items.Schema

	return nil
}

func (a *Array) Validate(ctx context.Context) error {
	return validation.ValidateObject(
		ctx,
		validation.ValidateField("items", a.Items),
		validation.ValidateField("minItems", validation.ValidatorFunc(func(ctx context.Context) error {
			if a.MinItems < 0 {
				return errors.New("must be greater than or equal to 0")
			}
			return nil
		})),
		validation.ValidatorFunc(func(ctx context.Context) error {
			if a.MaxItems != nil && a.MinItems > *a.MaxItems {
				return errors.New("minItems and maxItems constraints are impossible to fulfil")
			}
			return nil
		}),
		validateCommonFields(a, a.Const, a.Default, a.Enum),
	)
}

func (a *Array) ValidateSchema(s Schema, compare bool) error {
	if s.Type() == TypeNull {
		return nil
	}

	a2, ok := s.(*Array)
	if !ok {
		return typeErrorf("must be %s", a.TypeString())
	}

	if err := a.Items.ValidateSchema(a2.Items, compare); err != nil {
		if _, ok := err.(typeError); ok {
			return typeErrorf("must be %s", a.TypeString())
		}
		return err
	}

	return nil
}

func (a *Array) ValidateValue(value interface{}) (interface{}, error) {
	var v []interface{}
	if value != nil {
		var ok bool
		if v, ok = value.([]interface{}); !ok {
			return nil, errors.New("must be array")
		}
	}

	if a.Const != nil && a.CompareValues(a.Const, v) != 0 {
		return nil, fmt.Errorf("must be %s", a.StringValue(a.Const))
	}

	if len(a.Enum) == 1 && a.CompareValues(a.Enum[0], v) != 0 {
		return nil, fmt.Errorf("must be %s", a.StringValue(a.Enum[0]))
	}

	if len(a.Enum) > 0 {
		allowed := func() bool {
			for _, e := range a.Enum {
				if a.CompareValues(e, v) == 0 {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", a.join(a.Enum, ", "))
		}
	}

	ve := validation.Error{}

	if a.MaxItems != nil {
		if a.MinItems == *a.MaxItems && len(v) != int(a.MinItems) {
			switch a.MinItems {
			case 0:
				ve.AddError(errors.New("must be empty"))
			case 1:
				ve.AddError(errors.New("must have exactly one element"))
			default:
				ve.AddError(fmt.Errorf("must have exactly %d elements", a.MinItems))
			}
		} else {
			if len(v) > int(*a.MaxItems) {
				switch *a.MaxItems {
				case 0:
					ve.AddError(errors.New("must be empty"))
				case 1:
					ve.AddError(errors.New("must not contain more than one element"))
				default:
					ve.AddError(fmt.Errorf("must not contain more than %d elements", *a.MaxItems))
				}
			}
		}
	}

	if len(v) < int(a.MinItems) {
		switch a.MinItems {
		case 1:
			ve.AddError(errors.New("must have at least one element"))
		default:
			ve.AddError(fmt.Errorf("must have at least %d elements", a.MinItems))
		}
	}

	unique := true
	if a.UniqueItems {
		l := len(v)
		for i := 0; i < l; i++ {
			for j := i + 1; j < l; j++ {
				if a.Items.CompareValues(v[i], v[j]) == 0 {
					unique = false
					break
				}
			}
			if !unique {
				break
			}
		}
	}

	if !unique {
		ve.AddError(fmt.Errorf("array must contain unique items"))
	}

	for i, e := range v {
		nv, err := a.Items.ValidateValue(e)
		if err != nil {
			ve.AddFieldError(strconv.Itoa(i), err)
		} else {
			v[i] = nv
		}

	}

	if err := ve.ErrOrNil(); err != nil {
		return nil, err
	}

	return v, nil
}

func (a *Array) join(elems [][]interface{}, sep string) string {
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
