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

type Untyped struct {
	Metadata
	Types []string `json:"type,omitempty"`
}

func (u *Untyped) AssignValue(_ map[string]string, valueName, resultName string) string {
	if u.Pointer {
		panic("an untyped value can not have a pointer")
	}

	return fmt.Sprintf("%s = %s", resultName, valueName)
}

func (u *Untyped) CompareValues(v1, v2 interface{}) int {
	switch {
	case reflect.DeepEqual(v1, v2):
		return 0
	default:
		return -1
	}
}

func (u *Untyped) Copy() Schema {
	j, err := json.Marshal(u)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &Untyped{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (u *Untyped) DefaultValue() interface{} {
	return nil
}

func (u *Untyped) GoString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.Untyped{\n")
	if !reflect.ValueOf(u.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(u.Metadata.GoString()))
	}
	if len(u.Types) > 0 {
		_, _ = fmt.Fprintf(buf, "\tTypes: %#v,\n", u.Types)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (u *Untyped) GoType(_ map[string]string) string {
	if u.Pointer {
		panic("an untyped value can not have a pointer")
	}
	return "interface{}"
}

func (u *Untyped) StringValue(value interface{}) string {
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
		return strconv.Quote(v)
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}

		var sb strings.Builder
		sb.WriteRune('[')
		sb.WriteString(u.StringValue(v[0]))
		for _, e := range v[1:] {
			sb.WriteString(", ")
			sb.WriteString(u.StringValue(e))
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
		sb.WriteString(u.StringValue(v[keys[0]]))
		for _, k := range keys[1:] {
			sb.WriteString(", ")
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(u.StringValue(v[k]))
		}
		sb.WriteRune('}')
		return u.GoString()
	case io.Reader:
		return "<byte stream>"
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("<%T>", v)
	}
}

func (u *Untyped) Type() Type {
	return TypeUntyped
}

func (u *Untyped) TypeString() string {
	return "untyped"
}

func (u *Untyped) ValidateSchema(s Schema, compare bool) error {
	if len(u.Types) == 0 {
		return nil
	}

	isValid := false
	for _, t := range u.Types {
		if err := typeSchemas[Type(t)].ValidateSchema(s, compare); err == nil {
			isValid = true
			break
		}
	}

	if !isValid {
		if len(u.Types) == 1 {
			return fmt.Errorf("was expecting %s", u.Types[0])
		}
		return fmt.Errorf(
			"was expecting %s or %s",
			strings.Join(u.Types[0:len(u.Types)-1], ", "),
			u.Types[len(u.Types)-1],
		)
	}

	return nil
}

func (u *Untyped) ValidateValue(v interface{}) error {
	if len(u.Types) == 0 {
		return nil
	}

	isValid := false
	for _, t := range u.Types {
		if err := typeSchemas[Type(t)].ValidateValue(v); err == nil {
			isValid = true
			break
		}
	}

	if !isValid {
		if len(u.Types) == 1 {
			return fmt.Errorf("was expecting %s", u.Types[0])
		}
		return fmt.Errorf(
			"was expecting %s or %s",
			strings.Join(u.Types[0:len(u.Types)-1], ", "),
			u.Types[len(u.Types)-1],
		)
	}

	return nil
}

func UntypedValue() Schema {
	return untypedValueInst
}

var untypedValueInst = &untypedValue{
	Untyped: &Untyped{},
}

type untypedValue struct {
	*Untyped
}

func (u *untypedValue) Copy() Schema {
	return untypedValueInst
}

func (u *untypedValue) GoString() string {
	return "schema.UntypedValue()"
}
