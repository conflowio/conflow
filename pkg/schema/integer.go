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
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/conflowio/conflow/pkg/util/ptr"
	"github.com/conflowio/conflow/pkg/util/validation"
)

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type Integer struct {
	Metadata

	Const            *int64  `json:"const,omitempty"`
	Default          *int64  `json:"default,omitempty"`
	Enum             []int64 `json:"enum,omitempty"`
	ExclusiveMinimum *int64  `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *int64  `json:"exclusiveMaximum,omitempty"`
	// @enum ["int32", "int64"]
	Format     string `json:"format,omitempty"`
	Maximum    *int64 `json:"maximum,omitempty"`
	Minimum    *int64 `json:"minimum,omitempty"`
	Nullable   bool   `json:"nullable,omitempty"`
	MultipleOf *int64 `json:"multipleOf,omitempty"`
}

func (i *Integer) AssignValue(imports map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s(%s)", resultName, assignFuncName(i, imports), valueName)
}

func (i *Integer) CompareValues(v1, v2 interface{}) int {
	i1, _ := i.valueOf(v1)
	i2, _ := i.valueOf(v2)

	switch {
	case i1 == nil && i2 == nil:
		return 0
	case i1 == nil:
		return -1
	case i2 == nil:
		return 1
	case *i1 == *i2:
		return 0
	case *i1 < *i2:
		return -1
	default:
		return 1
	}
}

func (i *Integer) Copy() Schema {
	j, err := json.Marshal(i)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Integer{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (i *Integer) DefaultValue() interface{} {
	if i.Default == nil {
		return nil
	}
	return *i.Default
}

func (i *Integer) GetNullable() bool {
	return i.Nullable
}

func (i *Integer) GoString(imports map[string]string) string {
	pkg := schemaPkg(imports)
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sInteger{\n", pkg)
	if !reflect.ValueOf(i.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(i.Metadata.GoString(imports)))
	}
	if i.Const != nil {
		fprintf(buf, "\tConst: %sPointer(int64(%#v)),\n", pkg, *i.Const)
	}
	if i.Default != nil {
		fprintf(buf, "\tDefault: %sPointer(int64(%#v)),\n", pkg, *i.Default)
	}
	if len(i.Enum) > 0 {
		fprintf(buf, "\tEnum: %#v,\n", i.Enum)
	}
	if len(i.Format) > 0 {
		fprintf(buf, "\tFormat: %#v,\n", i.Format)
	}
	if i.Minimum != nil {
		fprintf(buf, "\tMinimum: %sPointer(int64(%#v)),\n", pkg, *i.Minimum)
	}
	if i.Maximum != nil {
		fprintf(buf, "\tMaximum: %sPointer(int64(%#v)),\n", pkg, *i.Maximum)
	}
	if i.ExclusiveMinimum != nil {
		fprintf(buf, "\tExclusiveMinimum: %sPointer(int64(%#v)),\n", pkg, *i.ExclusiveMinimum)
	}
	if i.ExclusiveMaximum != nil {
		fprintf(buf, "\tExclusiveMaximum: %sPointer(int64(%#v)),\n", pkg, *i.ExclusiveMaximum)
	}
	if i.MultipleOf != nil {
		fprintf(buf, "\tMultipleOf: %sPointer(int64(%#v)),\n", pkg, *i.MultipleOf)
	}
	if i.Nullable {
		fprintf(buf, "\tNullable: %#v,\n", i.Nullable)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (i *Integer) GoType(_ map[string]string) string {
	if i.Nullable {
		return "*int64"
	}
	return "int64"
}

func (i *Integer) MarshalJSON() ([]byte, error) {
	type Alias Integer
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(i.Type()),
		Alias: (*Alias)(i),
	})
}

func (i *Integer) SetNullable(nullable bool) {
	i.Nullable = nullable
}

func (i *Integer) StringValue(value interface{}) string {
	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		panic(fmt.Errorf("unexpected type when converting number to string: %T", v))
	}
}

func (i *Integer) Type() Type {
	return TypeInteger
}

func (i *Integer) TypeString() string {
	return string(TypeInteger)
}

func (i *Integer) UnmarshalJSON(input []byte) error {
	type Alias Integer
	return json.Unmarshal(input, &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(i),
	})
}

func (i *Integer) Validate(ctx *Context) error {
	return ValidateAll(ctx,
		func(ctx *Context) error {
			if i.ExclusiveMinimum != nil && i.Minimum != nil {
				return validation.NewFieldError("minimum", errors.New("should not be defined if exclusiveMinimum is set"))
			}

			if i.ExclusiveMaximum != nil && i.Maximum != nil {
				return validation.NewFieldError("maximum", errors.New("should not be defined if exclusiveMaximum is set"))
			}

			min := i.Minimum
			if i.ExclusiveMinimum != nil {
				min = ptr.To(*i.ExclusiveMinimum + 1)
			}

			max := i.Maximum
			if i.ExclusiveMaximum != nil {
				max = ptr.To(*i.ExclusiveMaximum - 1)
			}

			if min != nil && max != nil && *min > *max {
				return errors.New("minimum and maximum constraints are impossible to fulfil")
			}

			return nil
		},
		validateCommonFields(i, i.Const, i.Default, i.Enum),
	)
}

func (i *Integer) ValidateSchema(i2 Schema, compare bool) error {
	if compare {
		if i2.Type() != TypeInteger && i2.Type() != TypeNumber {
			return typeError("must be number")
		}
		return nil
	}

	if i2.Type() != TypeInteger {
		return typeError("must be integer")
	}

	return nil
}

func (i *Integer) ValidateValue(value interface{}) (interface{}, error) {
	v, ok := i.valueOf(value)
	if !ok {
		return nil, errors.New("must be integer")
	}
	if v == nil {
		return nil, nil
	}

	if i.Const != nil && *i.Const != *v {
		return nil, fmt.Errorf("must be %s", i.StringValue(*i.Const))
	}

	if len(i.Enum) == 1 && i.Enum[0] != *v {
		return nil, fmt.Errorf("must be %s", i.StringValue(i.Enum[0]))
	}

	if len(i.Enum) > 0 {
		allowed := func() bool {
			for _, e := range i.Enum {
				if e == *v {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", i.join(i.Enum, ", "))
		}
	}

	ve := &validation.Error{}

	if i.Format == "int32" && *v < math.MinInt32 {
		ve.AddErrorf("must be greater than or equal to %d", math.MinInt32)
	} else if i.ExclusiveMinimum != nil && *v <= *i.ExclusiveMinimum {
		ve.AddErrorf("must be greater than %d", *i.ExclusiveMinimum)
	} else if i.Minimum != nil && *v < *i.Minimum {
		ve.AddErrorf("must be greater than or equal to %d", *i.Minimum)
	}

	if i.Format == "int32" && *v > math.MaxInt32 {
		ve.AddErrorf("must be less than or equal to %d", math.MaxInt32)
	} else if i.ExclusiveMaximum != nil && *v >= *i.ExclusiveMaximum {
		ve.AddErrorf("must be less than %d", *i.ExclusiveMaximum)
	} else if i.Maximum != nil && *v > *i.Maximum {
		ve.AddErrorf("must be less than or equal to %d", *i.Maximum)
	}

	if i.MultipleOf != nil && *v%*i.MultipleOf != 0 {
		ve.AddErrorf("must be multiple of %d", *i.MultipleOf)
	}

	if err := ve.ErrOrNil(); err != nil {
		return nil, err
	}

	if i.Nullable {
		return v, nil
	}
	return *v, nil
}

func (i *Integer) join(elems []int64, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return i.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(i.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(i.StringValue(e))
	}
	return b.String()
}

func (i *Integer) valueOf(value interface{}) (*int64, bool) {
	switch v := value.(type) {
	case int64:
		return &v, true
	case *int64:
		return v, true
	default:
		return nil, false
	}
}

func IntegerValue() Schema {
	return integerValueInst
}

var integerValueInst = &integerValue{
	Integer: &Integer{},
}

type integerValue struct {
	*Integer
}

func (i *integerValue) Copy() Schema {
	return integerValueInst
}

func (i *integerValue) GoString(imports map[string]string) string {
	return fmt.Sprintf("%sIntegerValue()", schemaPkg(imports))
}
