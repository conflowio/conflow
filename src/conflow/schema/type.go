// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

type Type string

const (
	TypeUntyped    Type = ""
	TypeArray      Type = "array"
	TypeByteStream Type = "byte_stream"
	TypeBoolean    Type = "boolean"
	TypeFalse      Type = "false"
	TypeFunction   Type = "function"
	TypeInteger    Type = "integer"
	TypeMap        Type = "map"
	TypeNull       Type = "null"
	TypeNumber     Type = "number"
	TypeObject     Type = "object"
	TypeReference  Type = "reference"
	TypeString     Type = "string"
)

var typeSchemas = map[Type]Schema{
	TypeUntyped:    UntypedValue(),
	TypeArray:      &Array{Items: UntypedValue()},
	TypeByteStream: ByteStreamValue(),
	TypeBoolean:    BooleanValue(),
	TypeFalse:      False(),
	TypeFunction:   &Function{},
	TypeInteger:    IntegerValue(),
	TypeMap:        &Map{AdditionalProperties: UntypedValue()},
	TypeNull:       NullValue(),
	TypeNumber:     NumberValue(),
	TypeObject:     &Object{},
	TypeString:     StringValue(),
}
