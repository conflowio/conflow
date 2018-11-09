package util

import (
	"fmt"
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// NodeValueFunctionNames contains the type parser functions for every variable type
var NodeValueFunctionNames = map[string]string{
	basil.TypeAny:          "NodeAnyValue",
	basil.TypeArray:        "NodeArrayValue",
	basil.TypeBool:         "NodeBoolValue",
	basil.TypeFloat:        "NodeFloatValue",
	basil.TypeIdentifier:   "NodeIdentifierValue",
	basil.TypeInteger:      "NodeIntegerValue",
	basil.TypeMap:          "NodeMapValue",
	basil.TypeNumber:       "NodeNumberValue",
	basil.TypeString:       "NodeStringValue",
	basil.TypeStringArray:  "NodeStringArrayValue",
	basil.TypeTimeDuration: "NodeTimeDurationValue",
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
		return nil, parsley.NewError(node.Pos(), basil.ErrExpectingAny)
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

	return nil, parsley.NewError(node.Pos(), basil.ErrExpectingArray)
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

	return false, parsley.NewError(node.Pos(), basil.ErrExpectingBool)
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

	return 0.0, parsley.NewError(node.Pos(), basil.ErrExpectingFloat)
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

	return "", parsley.NewError(node.Pos(), basil.ErrExpectingIdentifier)
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

	return 0, parsley.NewError(node.Pos(), basil.ErrExpectingInteger)
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

	return nil, parsley.NewError(node.Pos(), basil.ErrExpectingMap)
}

// NodeNumberValue returns with the number value of a node
func NodeNumberValue(node parsley.Node, ctx interface{}) (basil.Number, parsley.Error) {
	val, err := node.Value(ctx)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return 0, nil
	}

	switch res := val.(type) {
	case int64:
		return res, nil
	case float64:
		return res, nil
	}

	return nil, parsley.NewError(node.Pos(), basil.ErrExpectingNumber)
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

	return "", parsley.NewError(node.Pos(), basil.ErrExpectingString)
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
				return nil, parsley.NewError(node.Pos(), basil.ErrExpectingStringArray)
			}
		}
		return res, nil
	}

	return nil, parsley.NewError(node.Pos(), basil.ErrExpectingString)
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

	return 0, parsley.NewError(node.Pos(), basil.ErrExpectingTimeDuration)
}

// CheckNodeType checks the type of the node
func CheckNodeType(node parsley.Node, expectedType string) parsley.Error {
	if node.Type() == basil.TypeUnknown || node.Type() == basil.TypeAny || expectedType == basil.TypeAny {
		return nil
	}

	if expectedType == basil.TypeNumber && (node.Type() == basil.TypeInteger || node.Type() == basil.TypeFloat) {
		return nil
	}

	if expectedType == basil.TypeArray && node.Type() == basil.TypeStringArray {
		return nil
	}

	if expectedType == basil.TypeStringArray && node.Type() == basil.TypeArray {
		for _, child := range node.(parsley.NonTerminalNode).Children() {
			if err := CheckNodeType(child, basil.TypeString); err != nil {
				return err
			}
		}
		return nil
	}

	if node.Type() != expectedType {
		typeErr, ok := basil.VariableTypeErrors[expectedType]
		if !ok {
			panic(fmt.Sprintf("Unknown type: %s", expectedType))
		}
		return parsley.NewError(node.Pos(), typeErr)
	}

	return nil
}
