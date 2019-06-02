package variable

import "errors"

// ErrNotDefined is an error when a variable is not defined
var ErrNotDefined = errors.New("variable not defined")

// Errors when expecting a certain variable type
var (
	ErrExpectingAny          = errors.New("was expecting any valid type")
	ErrExpectingArray        = errors.New("was expecting array")
	ErrExpectingBasic        = errors.New("was expecting any basic type")
	ErrExpectingBool         = errors.New("was expecting boolean")
	ErrExpectingFloat        = errors.New("was expecting float")
	ErrExpectingIdentifier   = errors.New("was expecting identifier")
	ErrExpectingInteger      = errors.New("was expecting integer")
	ErrExpectingMap          = errors.New("was expecting map")
	ErrExpectingNumber       = errors.New("was expecting number")
	ErrExpectingString       = errors.New("was expecting string")
	ErrExpectingStringArray  = errors.New("was expecting string array")
	ErrExpectingTime         = errors.New("was expecting RFC3339 time")
	ErrExpectingTimeDuration = errors.New("was expecting time duration")
	ErrExpectingWithLength   = errors.New("was expecting string, array or map")
)

// TypeErrors contains the type errors for all variable types
var TypeErrors = map[string]error{
	TypeAny:          ErrExpectingAny,
	TypeArray:        ErrExpectingArray,
	TypeBasic:        ErrExpectingBasic,
	TypeBool:         ErrExpectingBool,
	TypeFloat:        ErrExpectingFloat,
	TypeIdentifier:   ErrExpectingIdentifier,
	TypeInteger:      ErrExpectingInteger,
	TypeMap:          ErrExpectingMap,
	TypeNumber:       ErrExpectingNumber,
	TypeString:       ErrExpectingString,
	TypeStringArray:  ErrExpectingStringArray,
	TypeTime:         ErrExpectingTime,
	TypeTimeDuration: ErrExpectingTimeDuration,
	TypeWithLength:   ErrExpectingWithLength,
}
