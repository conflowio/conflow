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
	"io"
	"reflect"

	"github.com/conflowio/conflow/pkg/internal/utils"
)

type ByteStream struct {
	Metadata
}

func (b *ByteStream) AssignValue(imports map[string]string, valueName, resultName string) string {
	return fmt.Sprintf("%s = %s.(%sReadCloser)", resultName, valueName, utils.EnsureUniqueGoPackageSelector(imports, "io"))
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

func (b *ByteStream) GoString(imports map[string]string) string {
	buf := bytes.NewBuffer(nil)
	fprintf(buf, "&%sByteStream{\n", schemaPkg(imports))
	if !reflect.ValueOf(b.Metadata).IsZero() {
		fprintf(buf, "\tMetadata: %s,\n", indent(b.Metadata.GoString(imports)))
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

func (b *ByteStream) Validate(context.Context) error {
	return nil
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
	v, ok := value.(io.ReadCloser)
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

func (u *byteStreamValue) GoString(imports map[string]string) string {
	return fmt.Sprintf("%sByteStreamValue()", schemaPkg(imports))
}
