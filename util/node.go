package util

import (
	"time"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/parsley"
)

// NodeValueFunctionNames contains the type parser functions for every variable type
var NodeValueFunctionNames = map[string]string{
	ocl.TypeAny:          "NodeAnyValue",
	ocl.TypeArray:        "NodeArrayValue",
	ocl.TypeBool:         "NodeBoolValue",
	ocl.TypeFloat:        "NodeFloatValue",
	ocl.TypeInteger:      "NodeIntegerValue",
	ocl.TypeMap:          "NodeMapValue",
	ocl.TypeString:       "NodeStringValue",
	ocl.TypeTimeDuration: "NodeTimeDurationValue",
}

// NodeAnyValue returns with the array value of a node
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
		return nil, parsley.NewError(node.Pos(), ocl.ErrExpectingAny)
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

	return nil, parsley.NewError(node.Pos(), ocl.ErrExpectingArray)
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

	return false, parsley.NewError(node.Pos(), ocl.ErrExpectingBool)
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

	return 0.0, parsley.NewError(node.Pos(), ocl.ErrExpectingFloat)
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

	return 0, parsley.NewError(node.Pos(), ocl.ErrExpectingInteger)
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

	return nil, parsley.NewError(node.Pos(), ocl.ErrExpectingMap)
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

	return "", parsley.NewError(node.Pos(), ocl.ErrExpectingString)
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

	return 0, parsley.NewError(node.Pos(), ocl.ErrExpectingTimeDuration)
}
