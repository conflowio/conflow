package basil

import (
	"errors"

	"github.com/opsidian/parsley/text/terminal"
)

// Variable types
const (
	TypeAny          = "interface{}"
	TypeArray        = "[]interface{}"
	TypeBool         = terminal.BoolType
	TypeFloat        = terminal.FloatType
	TypeIdentifier   = "basil.ID"
	TypeInteger      = terminal.IntegerType
	TypeMap          = "map[string]interface{}"
	TypeString       = terminal.StringType
	TypeTimeDuration = terminal.TimeDurationType
)

// VariableTypes contains valid variable types with descriptions
var VariableTypes = map[string]string{
	TypeAny:          "any valid type",
	TypeArray:        "array",
	TypeBool:         "boolean",
	TypeFloat:        "float",
	TypeIdentifier:   "identifier",
	TypeInteger:      "integer",
	TypeMap:          "map",
	TypeString:       "string",
	TypeTimeDuration: "time duration",
}

// Errors when expecting a certain variable type
var (
	ErrExpectingAny          = errors.New("was expecting any valid type")
	ErrExpectingArray        = errors.New("was expecting array")
	ErrExpectingBool         = errors.New("was expecting boolean")
	ErrExpectingFloat        = errors.New("was expecting float")
	ErrExpectingIdentifier   = errors.New("was expecting identifier")
	ErrExpectingInteger      = errors.New("was expecting integer")
	ErrExpectingMap          = errors.New("was expecting map")
	ErrExpectingString       = errors.New("was expecting string")
	ErrExpectingTimeDuration = errors.New("was expecting time duration")
)

// VariableTypeErrors contains the type errors for all variable types
var VariableTypeErrors = map[string]error{
	TypeAny:          ErrExpectingAny,
	TypeArray:        ErrExpectingArray,
	TypeBool:         ErrExpectingBool,
	TypeFloat:        ErrExpectingFloat,
	TypeIdentifier:   ErrExpectingIdentifier,
	TypeInteger:      ErrExpectingInteger,
	TypeMap:          ErrExpectingMap,
	TypeString:       ErrExpectingString,
	TypeTimeDuration: ErrExpectingTimeDuration,
}
