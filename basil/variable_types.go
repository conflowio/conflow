package basil

import "errors"

// Variable types
const (
	TypeAny          = "interface{}"
	TypeArray        = "[]interface{}"
	TypeBool         = "bool"
	TypeFloat        = "float64"
	TypeInteger      = "int64"
	TypeMap          = "map[string]interface{}"
	TypeString       = "string"
	TypeTimeDuration = "time.Duration"
)

// VariableTypes contains valid variable types with descriptions
var VariableTypes = map[string]string{
	TypeAny:          "any valid type",
	TypeArray:        "array",
	TypeBool:         "boolean",
	TypeFloat:        "float",
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
	TypeInteger:      ErrExpectingInteger,
	TypeMap:          ErrExpectingMap,
	TypeString:       ErrExpectingString,
	TypeTimeDuration: ErrExpectingTimeDuration,
}
