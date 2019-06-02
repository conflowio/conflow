package variable

import (
	"time"

	"github.com/opsidian/basil/basil"
)

// ValueFunctionNames contains the type parser functions for every variable type
var ValueFunctionNames = map[string]string{
	TypeAny:          "AnyValue",
	TypeArray:        "ArrayValue",
	TypeBasic:        "BasicValue",
	TypeBool:         "BoolValue",
	TypeFloat:        "FloatValue",
	TypeIdentifier:   "IdentifierValue",
	TypeInteger:      "IntegerValue",
	TypeMap:          "MapValue",
	TypeNumber:       "NumberValue",
	TypeString:       "StringValue",
	TypeStringArray:  "StringArrayValue",
	TypeTime:         "TimeValue",
	TypeTimeDuration: "TimeDurationValue",
	TypeWithLength:   "WithLengthValue",
}

// AnyValue returns with any valid value
func AnyValue(val interface{}) (interface{}, error) {
	return val, nil
}

// ArrayValue returns with the array value of the given interface{} value
func ArrayValue(val interface{}) ([]interface{}, error) {
	if val == nil {
		return nil, nil
	}

	if res, ok := val.([]interface{}); ok {
		return res, nil
	}

	return nil, ErrExpectingArray
}

// BasicValue returns with a basic value
func BasicValue(val interface{}) (*Basic, error) {
	if val == nil {
		return nil, nil
	}

	if !IsBasicType(val) {
		return nil, ErrExpectingBasic
	}

	return NewBasic(val), nil
}

// BoolValue returns with the boolean value of the given interface{} value
func BoolValue(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}

	if res, ok := val.(bool); ok {
		return res, nil
	}

	return false, ErrExpectingBool
}

// FloatValue returns with the float value of the given interface{} value
func FloatValue(val interface{}) (float64, error) {
	if val == nil {
		return 0.0, nil
	}

	if res, ok := val.(float64); ok {
		return res, nil
	}

	return 0.0, ErrExpectingFloat
}

// IdentifierValue returns with the identifier value of the given interface{} value
func IdentifierValue(val interface{}) (basil.ID, error) {
	if val == nil {
		return "", nil
	}

	if res, ok := val.(basil.ID); ok {
		return res, nil
	}

	return "", ErrExpectingIdentifier
}

// IntegerValue returns with the integer value of the given interface{} value
func IntegerValue(val interface{}) (int64, error) {
	if val == nil {
		return 0, nil
	}

	if res, ok := val.(int64); ok {
		return res, nil
	}

	return 0, ErrExpectingInteger
}

// MapValue returns with the map value of the given interface{} value
func MapValue(val interface{}) (map[string]interface{}, error) {
	if val == nil {
		return nil, nil
	}

	if res, ok := val.(map[string]interface{}); ok {
		return res, nil
	}

	return nil, ErrExpectingMap
}

// NumberValue returns with the number value of the given interface{} value
func NumberValue(val interface{}) (*Number, error) {
	if val == nil {
		return nil, nil
	}

	if !IsNumberType(val) {
		return nil, ErrExpectingNumber
	}

	return NewNumber(val), nil
}

// StringValue returns with the string value of the given interface{} value
func StringValue(val interface{}) (string, error) {
	if val == nil {
		return "", nil
	}

	if res, ok := val.(string); ok {
		return res, nil
	}

	return "", ErrExpectingString
}

// StringArrayValue returns with the string array value of the given interface{} value
func StringArrayValue(val interface{}) ([]string, error) {
	if val == nil {
		return nil, nil
	}

	switch v := val.(type) {
	case []string:
		return v, nil
	case []interface{}:
		var ok bool
		res := make([]string, len(v))
		for i := range v {
			if res[i], ok = v[i].(string); !ok {
				return nil, ErrExpectingStringArray
			}
		}
		return res, nil
	}

	return nil, ErrExpectingString
}

// TimeValue returns with the time  value of the given interface{} value
// It accepts a time.Time object or a time string in RFC3339 format
func TimeValue(val interface{}) (time.Time, error) {
	if val == nil {
		return time.Time{}, nil
	}

	if res, ok := val.(time.Time); ok {
		return res, nil
	}

	if res, ok := val.(string); ok {
		t, err := time.Parse(time.RFC3339, res)
		if err != nil {
			return time.Time{}, err
		}

		return t, nil
	}

	return time.Time{}, ErrExpectingTime
}

// TimeDurationValue returns with the time duration value of the given interface{} value
func TimeDurationValue(val interface{}) (time.Duration, error) {
	if val == nil {
		return 0, nil
	}

	if res, ok := val.(time.Duration); ok {
		return res, nil
	}

	return 0, ErrExpectingTimeDuration
}

// WithLengthValue returns with a value which has a length
func WithLengthValue(val interface{}) (*WithLength, error) {
	if val == nil {
		return nil, nil
	}

	if !IsWithLengthType(val) {
		return nil, ErrExpectingWithLength
	}

	return NewWithLength(val), nil
}
