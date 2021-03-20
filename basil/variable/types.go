// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package variable

import (
	"github.com/opsidian/parsley/text/terminal"
)

// Variable types
const (
	TypeAny          = "interface{}"
	TypeArray        = "[]interface{}"
	TypeBasic        = "variable.Basic"
	TypeBool         = terminal.BoolType
	TypeFloat        = terminal.FloatType
	TypeIdentifier   = "basil.ID"
	TypeInteger      = terminal.IntegerType
	TypeMap          = "map[string]interface{}"
	TypeNil          = "nil"
	TypeNumber       = "variable.Number"
	TypeStream       = "io.ReadCloser"
	TypeString       = terminal.StringType
	TypeStringArray  = "[]string"
	TypeTime         = "time.Time"
	TypeTimeDuration = terminal.TimeDurationType
	TypeWithLength   = "variable.WithLength"
	TypeUnknown      = ""
)

// Types contains valid variable types with descriptions
var Types = map[string]string{
	TypeAny:          "any valid type",
	TypeArray:        "array",
	TypeBasic:        "any basic type",
	TypeBool:         "boolean",
	TypeFloat:        "float",
	TypeIdentifier:   "identifier",
	TypeInteger:      "integer",
	TypeMap:          "map",
	TypeNil:          "nil",
	TypeNumber:       "number",
	TypeStream:       "stream",
	TypeString:       "string",
	TypeStringArray:  "string array",
	TypeTime:         "time",
	TypeTimeDuration: "time duration",
	TypeWithLength:   "string, array or map",
}

// UnionTypes contains all union variable types
var UnionTypes = map[string][]string{
	TypeArray:      []string{TypeStringArray},
	TypeBasic:      []string{TypeBool, TypeFloat, TypeIdentifier, TypeInteger, TypeNumber, TypeString, TypeTime, TypeTimeDuration},
	TypeNumber:     []string{TypeFloat, TypeInteger},
	TypeWithLength: []string{TypeArray, TypeIdentifier, TypeString, TypeStringArray, TypeMap},
}
