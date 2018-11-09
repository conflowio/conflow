package basil

import (
	"errors"
	"time"

	"github.com/opsidian/parsley/text/terminal"
)

// Number is an integer or float number
type Number interface{}

// Variable types
const (
	TypeAny          = "interface{}"
	TypeArray        = "[]interface{}"
	TypeBool         = terminal.BoolType
	TypeFloat        = terminal.FloatType
	TypeIdentifier   = "basil.ID"
	TypeInteger      = terminal.IntegerType
	TypeMap          = "map[string]interface{}"
	TypeNumber       = "basil.Number"
	TypeString       = terminal.StringType
	TypeStringArray  = "[]string"
	TypeTimeDuration = terminal.TimeDurationType
	TypeUnknown      = ""
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
	TypeNumber:       "number",
	TypeString:       "string",
	TypeStringArray:  "string array",
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
	ErrExpectingNumber       = errors.New("was expecting number")
	ErrExpectingString       = errors.New("was expecting string")
	ErrExpectingStringArray  = errors.New("was expecting string array")
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
	TypeNumber:       ErrExpectingNumber,
	TypeString:       ErrExpectingString,
	TypeStringArray:  ErrExpectingStringArray,
	TypeTimeDuration: ErrExpectingTimeDuration,
}

// GetValueType returns with the type of the given value
func GetValueType(value interface{}) string {
	switch value.(type) {
	case []interface{}:
		return TypeArray
	case bool:
		return TypeBool
	case float64:
		return TypeFloat
	case ID:
		return TypeIdentifier
	case int64:
		return TypeInteger
	case map[string]interface{}:
		return TypeMap
	case Number:
		return TypeNumber
	case string:
		return TypeString
	case []string:
		return TypeStringArray
	case time.Duration:
		return TypeTimeDuration
	default:
		return TypeUnknown
	}
}
