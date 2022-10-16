// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type Any struct {
	Metadata
	Types []string `json:"type,omitempty"`
}

func (a *Any) AssignValue(_ map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s", resultName, valueName)
}

func (a *Any) CompareValues(v1, v2 interface{}) int {
	switch {
	case reflect.DeepEqual(v1, v2):
		return 0
	default:
		return -1
	}
}

func (a *Any) Copy() Schema {
	j, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Any{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (a *Any) DefaultValue() interface{} {
	return nil
}

func (a *Any) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Any{\n")
	if !reflect.ValueOf(a.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(a.Metadata.GoString(imports)))
	}
	if len(a.Types) > 0 {
		_, _ = fmt.Fprintf(buf, "\tTypes: %#v,\n", a.Types)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (a *Any) GoType(_ map[string]string) string {
	return "interface{}"
}

func (a *Any) StringValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case bool:
		return strconv.FormatBool(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}

		var sb strings.Builder
		sb.WriteRune('[')
		sb.WriteString(a.StringValue(v[0]))
		for _, e := range v[1:] {
			sb.WriteString(", ")
			sb.WriteString(a.StringValue(e))
		}
		sb.WriteRune(']')
		return sb.String()
	case map[string]interface{}:
		if len(v) == 0 {
			return "{}"
		}

		keys := getSortedMapKeys(v)

		var sb strings.Builder
		sb.WriteRune('{')
		sb.WriteString(keys[0])
		sb.WriteString(": ")
		sb.WriteString(a.StringValue(v[keys[0]]))
		for _, k := range keys[1:] {
			sb.WriteString(", ")
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(a.StringValue(v[k]))
		}
		sb.WriteRune('}')
		return sb.String()
	case io.Reader:
		return "<byte stream>"
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("<%T>", v)
	}
}

func (a *Any) Type() Type {
	return TypeAny
}

func (a *Any) TypeString() string {
	return "any"
}

func (a *Any) Validate(*Context) error {
	return nil
}

func (a *Any) ValidateSchema(s Schema, compare bool) error {
	if len(a.Types) == 0 {
		return nil
	}

	isValid := false
	for _, t := range a.Types {
		if err := typeSchemas[Type(t)].ValidateSchema(s, compare); err == nil {
			isValid = true
			break
		}
	}

	if !isValid {
		if len(a.Types) == 1 {
			return fmt.Errorf("was expecting %s", a.Types[0])
		}
		return fmt.Errorf(
			"was expecting %s or %s",
			strings.Join(a.Types[0:len(a.Types)-1], ", "),
			a.Types[len(a.Types)-1],
		)
	}

	return nil
}

func (a *Any) ValidateValue(v interface{}) (interface{}, error) {
	if len(a.Types) == 0 {
		return v, nil
	}

	isValid := false
	for _, t := range a.Types {
		nv, err := typeSchemas[Type(t)].ValidateValue(v)
		if err == nil {
			v = nv
			isValid = true
			break
		}
	}

	if !isValid {
		if len(a.Types) == 1 {
			return nil, fmt.Errorf("was expecting %s", a.Types[0])
		}
		return nil, fmt.Errorf(
			"was expecting %s or %s",
			strings.Join(a.Types[0:len(a.Types)-1], ", "),
			a.Types[len(a.Types)-1],
		)
	}

	return v, nil
}

func AnyValue() Schema {
	return anyValueInst
}

var anyValueInst = &anyValue{
	Any: &Any{},
}

type anyValue struct {
	*Any
}

func (u *anyValue) Copy() Schema {
	return anyValueInst
}

func (u *anyValue) GoString(map[string]string) string {
	return "schema.AnyValue()"
}
