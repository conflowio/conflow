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
	"io"
	"reflect"

	"github.com/conflowio/conflow/src/internal/utils"
)

type ByteStream struct {
	Metadata
}

func (b *ByteStream) AssignValue(imports map[string]string, valueName, resultName string) string {
	ioPackageName := utils.EnsureUniqueGoPackageName(imports, "io")

	return fmt.Sprintf("%s = %s.(%s.ReadCloser)", resultName, valueName, ioPackageName)
}

func (b *ByteStream) CompareValues(v1, v2 interface{}) int {
	return -1
}

func (b *ByteStream) Copy() Schema {
	j, err := json.Marshal(b)
	if err != nil {
		panic(fmt.Errorf("failed to encode schema: %w", err))
	}

	cp := &ByteStream{}
	if err := json.Unmarshal(j, cp); err != nil {
		panic(fmt.Errorf("failed to decode schema: %w", err))
	}

	return cp
}

func (b *ByteStream) DefaultValue() interface{} {
	return nil
}

func (b *ByteStream) GoString(map[string]string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("&schema.ByteStream{\n")
	if !reflect.ValueOf(b.Metadata).IsZero() {
		_, _ = fmt.Fprintf(buf, "\tMetadata: %s,\n", indent(b.Metadata.GoString()))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (b *ByteStream) GoType(imports map[string]string) string {
	return utils.GoType(imports, reflect.TypeOf(io.ReadCloser(nil)), false)
}

func (b *ByteStream) MarshalJSON() ([]byte, error) {
	type Alias ByteStream
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(b.Type()),
		Alias: (*Alias)(b),
	})
}

func (b *ByteStream) StringValue(value interface{}) string {
	return "<byte stream>"
}

func (b *ByteStream) Type() Type {
	return TypeByteStream
}

func (b *ByteStream) TypeString() string {
	return string(TypeByteStream)
}

func (b *ByteStream) ValidateSchema(s2 Schema, compare bool) error {
	if compare {
		return fmt.Errorf("byte streams are not comparable")
	}

	if s2.Type() != TypeByteStream {
		return typeError("must be byte stream")
	}

	return nil
}

func (b *ByteStream) ValidateValue(value interface{}) (interface{}, error) {
	v, ok := value.(io.Reader)
	if !ok {
		return nil, errors.New("must be byte stream")
	}

	return v, nil
}

func ByteStreamValue() Schema {
	return byteStreamValueInst
}

var byteStreamValueInst = &byteStreamValue{
	ByteStream: &ByteStream{},
}

type byteStreamValue struct {
	*ByteStream
}

func (u *byteStreamValue) Copy() Schema {
	return byteStreamValueInst
}

func (u *byteStreamValue) GoString(map[string]string) string {
	return "schema.ByteStreamValue()"
}
