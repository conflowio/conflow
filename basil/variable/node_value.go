package variable

import (
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// NodeValueFunctionNames contains the type parser functions for every variable type
var NodeValueFunctionNames = map[string]string{
	TypeAny:          "NodeAnyValue",
	TypeArray:        "NodeArrayValue",
	TypeBasic:        "NodeBasicValue",
	TypeBool:         "NodeBoolValue",
	TypeFloat:        "NodeFloatValue",
	TypeIdentifier:   "NodeIdentifierValue",
	TypeInteger:      "NodeIntegerValue",
	TypeMap:          "NodeMapValue",
	TypeNumber:       "NodeNumberValue",
	TypeString:       "NodeStringValue",
	TypeStringArray:  "NodeStringArrayValue",
	TypeTimeDuration: "NodeTimeDurationValue",
	TypeWithLength:   "NodeWithLengthValue",
}

// NodeAnyValue returns with any valid value
func NodeAnyValue(node parsley.Node, ctx interface{}) (interface{}, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	switch val.(type) {
	case []interface{}:
	case bool:
	case float64:
	case int64:
	case map[string]interface{}:
	case string:
	case time.Duration:
	default:
		return nil, parsley.NewError(node.Pos(), ErrExpectingAny)
	}

	return val, nil
}

// NodeArrayValue returns with the array value of a node
func NodeArrayValue(node parsley.Node, ctx interface{}) ([]interface{}, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	if res, ok := val.([]interface{}); ok {
		return res, nil
	}

	return nil, parsley.NewError(node.Pos(), ErrExpectingArray)
}

// NodeBasicValue returns with a basic value
func NodeBasicValue(node parsley.Node, ctx interface{}) (*Basic, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	if !IsBasicType(val) {
		return nil, parsley.NewError(node.Pos(), ErrExpectingBasic)
	}

	return NewBasic(val), nil
}

// NodeBoolValue returns with the boolean value of a node
func NodeBoolValue(node parsley.Node, ctx interface{}) (bool, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return false, err
	}

	if val == nil {
		return false, nil
	}

	if res, ok := val.(bool); ok {
		return res, nil
	}

	return false, parsley.NewError(node.Pos(), ErrExpectingBool)
}

// NodeFloatValue returns with the float value of a node
func NodeFloatValue(node parsley.Node, ctx interface{}) (float64, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return 0.0, err
	}

	if val == nil {
		return 0.0, nil
	}

	if res, ok := val.(float64); ok {
		return res, nil
	}

	return 0.0, parsley.NewError(node.Pos(), ErrExpectingFloat)
}

// NodeIdentifierValue returns with the identifier value of a node
func NodeIdentifierValue(node parsley.Node, ctx interface{}) (basil.ID, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return "", err
	}

	if val == nil {
		return "", nil
	}

	if res, ok := val.(basil.ID); ok {
		return res, nil
	}

	return "", parsley.NewError(node.Pos(), ErrExpectingIdentifier)
}

// NodeIntegerValue returns with the integer value of a node
func NodeIntegerValue(node parsley.Node, ctx interface{}) (int64, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return 0, err
	}

	if val == nil {
		return 0, nil
	}

	if res, ok := val.(int64); ok {
		return res, nil
	}

	return 0, parsley.NewError(node.Pos(), ErrExpectingInteger)
}

// NodeMapValue returns with the map value of a node
func NodeMapValue(node parsley.Node, ctx interface{}) (map[string]interface{}, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	if res, ok := val.(map[string]interface{}); ok {
		return res, nil
	}

	return nil, parsley.NewError(node.Pos(), ErrExpectingMap)
}

// NodeNumberValue returns with the number value of a node
func NodeNumberValue(node parsley.Node, ctx interface{}) (*Number, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	if !IsNumberType(val) {
		return nil, parsley.NewError(node.Pos(), ErrExpectingNumber)
	}

	return NewNumber(val), nil
}

// NodeStringValue returns with the string value of a node
func NodeStringValue(node parsley.Node, ctx interface{}) (string, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return "", err
	}

	if val == nil {
		return "", nil
	}

	if res, ok := val.(string); ok {
		return res, nil
	}

	return "", parsley.NewError(node.Pos(), ErrExpectingString)
}

// NodeStringArrayValue returns with the string array value of a node
func NodeStringArrayValue(node parsley.Node, ctx interface{}) ([]string, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

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
				return nil, parsley.NewError(node.Pos(), ErrExpectingStringArray)
			}
		}
		return res, nil
	}

	return nil, parsley.NewError(node.Pos(), ErrExpectingString)
}

// NodeTimeDurationValue returns with the time duration value of a node
func NodeTimeDurationValue(node parsley.Node, ctx interface{}) (time.Duration, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return 0, err
	}

	if val == nil {
		return 0, nil
	}

	if res, ok := val.(time.Duration); ok {
		return res, nil
	}

	return 0, parsley.NewError(node.Pos(), ErrExpectingTimeDuration)
}

// NodeWithLengthValue returns with a value which has a length
func NodeWithLengthValue(node parsley.Node, ctx interface{}) (*WithLength, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	if !IsWithLengthType(val) {
		return nil, parsley.NewError(node.Pos(), ErrExpectingWithLength)
	}

	return NewWithLength(val), nil
}
