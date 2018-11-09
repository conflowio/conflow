package variable

import "github.com/opsidian/parsley/text/terminal"

// Variable types
const (
	TypeAny          = "interface{}"
	TypeArray        = "[]interface{}"
	TypeBasic        = "*variable.Basic"
	TypeBool         = terminal.BoolType
	TypeFloat        = terminal.FloatType
	TypeIdentifier   = "variable.ID"
	TypeInteger      = terminal.IntegerType
	TypeMap          = "map[string]interface{}"
	TypeNumber       = "*variable.Number"
	TypeString       = terminal.StringType
	TypeStringArray  = "[]string"
	TypeTimeDuration = terminal.TimeDurationType
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
	TypeNumber:       "number",
	TypeString:       "string",
	TypeStringArray:  "string array",
	TypeTimeDuration: "time duration",
}

// UnionTypes contains all union variable types
var UnionTypes = map[string][]string{
	TypeBasic:  []string{TypeBool, TypeFloat, TypeIdentifier, TypeInteger, TypeNumber, TypeString, TypeTimeDuration},
	TypeNumber: []string{TypeFloat, TypeInteger},
	TypeArray:  []string{TypeStringArray},
}
