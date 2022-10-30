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
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/conflowio/conflow/pkg/util/ptr"
	"github.com/conflowio/conflow/pkg/util/validation"
)

// Epsilon is used as a float64 comparison tolerance
const Epsilon = 0.000000001

//	@block {
//	  type = "configuration"
//	  path = "interpreters"
//	}
type Number struct {
	Metadata

	Const            *float64  `json:"const,omitempty"`
	Default          *float64  `json:"default,omitempty"`
	Enum             []float64 `json:"enum,omitempty"`
	ExclusiveMinimum *float64  `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *float64  `json:"exclusiveMaximum,omitempty"`
	Maximum          *float64  `json:"maximum,omitempty"`
	Minimum          *float64  `json:"minimum,omitempty"`
	MultipleOf       *float64  `json:"multipleOf,omitempty"`
	Nullable         bool      `json:"nullable,omitempty"`
}

func (n *Number) AssignValue(imports map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s(%s)", resultName, assignFuncName(n, imports), valueName)
}

func (n *Number) CompareValues(v1, v2 interface{}) int {
	n1, _ := n.valueOf(v1)
	n2, _ := n.valueOf(v2)

	switch {
	case n1 == nil && n2 == nil:
		return 0
	case n1 == nil:
		return -1
	case n2 == nil:
		return 1
	case *n1-*n2 < Epsilon && *n2-*n1 < Epsilon:
		return 0
	case *n1 < *n2:
		return -1
	default:
		return 1
	}
}

func (n *Number) Copy() Schema {
	j, err := json.Marshal(n)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Number{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (n *Number) DefaultValue() interface{} {
	if n.Default == nil {
		return nil
	}
	return *n.Default
}

func (n *Number) MarshalJSON() ([]byte, error) {
	type Alias Number
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(n.Type()),
		Alias: (*Alias)(n),
	})
}

func (n *Number) GetNullable() bool {
	return n.Nullable
}

func (n *Number) GoString(imports map[string]string) string {
	pkg := schemaPkg(imports)
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sNumber{\n", pkg)
	if !reflect.ValueOf(n.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(n.Metadata.GoString(imports)))
	}
	if n.Const != nil {
		fprintf(buf, "\tConst: %sPointer(float64(%#v)),\n", pkg, *n.Const)
	}
	if n.Default != nil {
		fprintf(buf, "\tDefault: %sPointer(float64(%#v)),\n", pkg, *n.Default)
	}
	if len(n.Enum) > 0 {
		fprintf(buf, "\tEnum: %#v,\n", n.Enum)
	}
	if n.Minimum != nil {
		fprintf(buf, "\tMinimum: %sPointer(float64(%#v)),\n", pkg, *n.Minimum)
	}
	if n.Maximum != nil {
		fprintf(buf, "\tMaximum: %sPointer(float64(%#v)),\n", pkg, *n.Maximum)
	}
	if n.ExclusiveMinimum != nil {
		fprintf(buf, "\tExclusiveMinimum: %sPointer(float64(%#v)),\n", pkg, *n.ExclusiveMinimum)
	}
	if n.ExclusiveMaximum != nil {
		fprintf(buf, "\tExclusiveMaximum: %sPointer(float64(%#v)),\n", pkg, *n.ExclusiveMaximum)
	}
	if n.MultipleOf != nil {
		fprintf(buf, "\tMultipleOf: %sPointer(float64(%#v)),\n", pkg, *n.MultipleOf)
	}
	if n.Nullable {
		fprintf(buf, "\tNullable: %#v,\n", n.Nullable)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (n *Number) GoType(_ map[string]string) string {
	if n.Nullable {
		return "*float64"
	}
	return "float64"
}

func (n *Number) SetNullable(nullable bool) {
	n.Nullable = nullable
}

func (n *Number) StringValue(value interface{}) string {
	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func (n *Number) Type() Type {
	return TypeNumber
}

func (n *Number) TypeString() string {
	return string(TypeNumber)
}

func (n *Number) UnmarshalJSON(input []byte) error {
	type Alias Number
	return json.Unmarshal(input, &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}

func (n *Number) Validate(ctx context.Context) error {
	if n.ExclusiveMinimum != nil && n.Minimum != nil {
		return validation.NewFieldError("minimum", errors.New("should not be defined if exclusiveMinimum is set"))
	}

	if n.ExclusiveMaximum != nil && n.Maximum != nil {
		return validation.NewFieldError("maximum", errors.New("should not be defined if exclusiveMaximum is set"))
	}

	min := n.Minimum
	if n.ExclusiveMinimum != nil {
		min = ptr.To(*n.ExclusiveMinimum + Epsilon)
	}

	max := n.Maximum
	if n.ExclusiveMaximum != nil {
		max = ptr.To(*n.ExclusiveMaximum - Epsilon)
	}

	if min != nil && max != nil && NumberGreaterThan(*min, *max) {
		return errors.New("minimum and maximum constraints are impossible to fulfil")
	}

	if err := validateCommonFields(n, n.Const, n.Default, n.Enum); err != nil {
		return err
	}

	return nil
}

func (n *Number) ValidateSchema(n2 Schema, _ bool) error {
	if n2.Type() != TypeNumber && n2.Type() != TypeInteger {
		return typeError("must be number")
	}

	return nil
}

func (n *Number) ValidateValue(value interface{}) (interface{}, error) {
	v, ok := n.valueOf(value)
	if !ok {
		return nil, errors.New("must be number")
	}
	if v == nil {
		return nil, nil
	}

	if n.Const != nil && !NumberEquals(*n.Const, *v) {
		return nil, fmt.Errorf("must be %s", n.StringValue(*n.Const))
	}

	if len(n.Enum) == 1 && !NumberEquals(n.Enum[0], *v) {
		return nil, fmt.Errorf("must be %s", n.StringValue(n.Enum[0]))
	}

	if len(n.Enum) > 0 {
		allowed := func() bool {
			for _, e := range n.Enum {
				if NumberEquals(e, *v) {
					return true
				}
			}
			return false
		}
		if !allowed() {
			return nil, fmt.Errorf("must be one of %s", n.join(n.Enum, ", "))
		}
	}

	ve := &validation.Error{}

	if n.ExclusiveMinimum != nil && !NumberGreaterThan(*v, *n.ExclusiveMinimum) {
		ve.AddErrorf("must be greater than %s", n.StringValue(*n.ExclusiveMinimum))
	} else if n.Minimum != nil && !NumberGreaterThanOrEqualsTo(*v, *n.Minimum) {
		ve.AddErrorf("must be greater than or equal to %s", n.StringValue(*n.Minimum))
	}

	if n.ExclusiveMaximum != nil && !NumberLessThan(*v, *n.ExclusiveMaximum) {
		ve.AddErrorf("must be less than %s", n.StringValue(*n.ExclusiveMaximum))
	} else if n.Maximum != nil && !NumberLessThanOrEqualsTo(*v, *n.Maximum) {
		ve.AddErrorf("must be less than or equal to %s", n.StringValue(*n.Maximum))
	}

	if n.MultipleOf != nil && !NumberMultipleOf(*v, *n.MultipleOf) {
		ve.AddErrorf("must be multiple of %s", n.StringValue(*n.MultipleOf))
	}

	if err := ve.ErrOrNil(); err != nil {
		return nil, err
	}

	if n.Nullable {
		return v, nil
	}
	return *v, nil
}

func (n *Number) join(elems []float64, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return n.StringValue(elems[0])
	}

	var b strings.Builder
	b.WriteString(n.StringValue(elems[0]))
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(n.StringValue(e))
	}
	return b.String()
}

func (n *Number) valueOf(value interface{}) (*float64, bool) {
	switch v := value.(type) {
	case int64:
		return ptr.To(float64(v)), true
	case *int64:
		return ptr.To(float64(*v)), true
	case float64:
		return &v, true
	case *float64:
		return v, true
	default:
		return nil, false
	}
}

func NumberValue() Schema {
	return numberValueInst
}

var numberValueInst = &numberValue{
	Number: &Number{},
}

type numberValue struct {
	*Number
}

func (n *numberValue) Copy() Schema {
	return numberValueInst
}

func (n *numberValue) GoString(imports map[string]string) string {
	return fmt.Sprintf("%sNumberValue()", schemaPkg(imports))
}

func NumberEquals(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

func NumberLessThan(v1, v2 float64) bool {
	return v1 <= v2-Epsilon
}

func NumberLessThanOrEqualsTo(v1, v2 float64) bool {
	return v1 < v2+Epsilon
}

func NumberGreaterThan(v1, v2 float64) bool {
	return v1 >= v2+Epsilon
}

func NumberGreaterThanOrEqualsTo(v1, v2 float64) bool {
	return v1 > v2-Epsilon
}

func NumberMultipleOf(v1, v2 float64) bool {
	div := v1 / v2
	return math.Abs(div-math.Round(div)) < Epsilon
}
