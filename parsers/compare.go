// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"
	"reflect"
	"time"

	"github.com/opsidian/basil/basil/schema"

	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// Compare will match comparison expressions defined by the following rule, where P is the input parser:
//  S       -> P (COMP_OP P)*
//  COMP_OP -> "=="
//          -> "!="
//          -> "<"
//          -> "<="
//          -> ">"
//          -> ">="
func Compare(p parsley.Parser) *combinator.Sequence {
	return SepByOp(
		p,
		combinator.Choice(
			terminal.Op("=="),
			terminal.Op("!="),
			terminal.Op("<="),
			terminal.Op("<"),
			terminal.Op(">="),
			terminal.Op(">"),
		),
	).Token("COMPARE").Bind(compareInterpreter{}).ReturnSingle()
}

type compareInterpreter struct{}

func (a compareInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	var resultSchema schema.Schema
	var op string
	var opPos parsley.Pos
	expectsOp := false

	for i, node := range node.Children() {
		if i == 0 {
			resultSchema = node.Schema().(schema.Schema)
		} else if expectsOp {
			op = node.(parsley.LiteralNode).Value().(string)
			opPos = node.Pos()
		} else {
			if err := resultSchema.ValidateSchema(node.Schema().(schema.Schema), true); err != nil {
				return nil, parsley.NewErrorf(
					opPos,
					"unsupported %s operation on %s and %s",
					op,
					resultSchema.TypeString(),
					node.Schema().(schema.Schema).TypeString(),
				)
			}

			if op != "==" && op != "!=" {
				switch resultSchema.(type) {
				case schema.ArrayKind:
					return "", parsley.NewErrorf(opPos, "%q operator is invalid for arrays", op)
				case schema.MapKind:
					return "", parsley.NewErrorf(opPos, "%q operator is invalid for maps", op)
				case schema.ObjectKind:
					return "", parsley.NewErrorf(opPos, "%q operator is invalid for objects", op)
				}
			}

			resultSchema = schema.BooleanValue()
		}
		expectsOp = !expectsOp
	}
	return resultSchema, nil
}

func (a compareInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	var res interface{}
	var op string
	var opPos parsley.Pos
	var err parsley.Error
	expectsOp := false

	for i, node := range nodes {
		var v interface{}
		v, err = parsley.EvaluateNode(ctx, node)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			res = v
		} else if expectsOp {
			op = v.(string)
			opPos = node.Pos()
		} else {
			switch vt := v.(type) {
			case bool:
				res, err = compareBool(res, op, vt, opPos)
			case int64:
				res, err = compareInt(res, op, vt, opPos)
			case float64:
				res, err = compareFloat(res, op, vt, opPos)
			case string:
				res, err = compareString(res, op, vt, opPos)
			case time.Time:
				res, err = compareTime(res, op, vt, opPos)
			case time.Duration:
				res, err = compareTimeDuration(res, op, vt, opPos)
			default:
				res, err = compareOther(res, op, v, opPos)
			}
		}
		if err != nil {
			return nil, err
		}
		expectsOp = !expectsOp
	}
	return res, nil
}

func floatsEqual(f1 float64, f2 float64) bool {
	return f1-f2 < schema.Epsilon && f2-f1 < schema.Epsilon
}

func compareInt(res interface{}, op string, v int64, opPos parsley.Pos) (interface{}, parsley.Error) {
	switch op {
	case "==":
		switch rest := res.(type) {
		case int64:
			return rest == v, nil
		case float64:
			return floatsEqual(rest, float64(v)), nil
		}
	case "!=":
		switch rest := res.(type) {
		case int64:
			return rest != v, nil
		case float64:
			return !floatsEqual(rest, float64(v)), nil
		}
	case "<":
		switch rest := res.(type) {
		case int64:
			return rest < v, nil
		case float64:
			return !floatsEqual(rest, float64(v)) && rest < float64(v), nil
		}
	case "<=":
		switch rest := res.(type) {
		case int64:
			return rest <= v, nil
		case float64:
			return floatsEqual(rest, float64(v)) || rest < float64(v), nil
		}
	case ">":
		switch rest := res.(type) {
		case int64:
			return rest > v, nil
		case float64:
			return !floatsEqual(rest, float64(v)) && rest > float64(v), nil
		}
	case ">=":
		switch rest := res.(type) {
		case int64:
			return rest >= v, nil
		case float64:
			return floatsEqual(rest, float64(v)) || rest > float64(v), nil
		}
	}
	return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", op, fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
}

func compareFloat(res interface{}, op string, v float64, opPos parsley.Pos) (interface{}, parsley.Error) {
	switch op {
	case "==":
		switch rest := res.(type) {
		case int64:
			return floatsEqual(float64(rest), v), nil
		case float64:
			return floatsEqual(rest, v), nil
		}
	case "!=":
		switch rest := res.(type) {
		case int64:
			return !floatsEqual(float64(rest), v), nil
		case float64:
			return !floatsEqual(rest, v), nil
		}
	case "<":
		switch rest := res.(type) {
		case int64:
			return !floatsEqual(float64(rest), v) && float64(rest) < v, nil
		case float64:
			return !floatsEqual(rest, v) && rest < v, nil
		}
	case "<=":
		switch rest := res.(type) {
		case int64:
			return floatsEqual(float64(rest), v) || float64(rest) < v, nil
		case float64:
			return floatsEqual(rest, v) || rest < v, nil
		}
	case ">":
		switch rest := res.(type) {
		case int64:
			return !floatsEqual(float64(rest), v) && float64(rest) > v, nil
		case float64:
			return !floatsEqual(rest, v) && rest > v, nil
		}
	case ">=":
		switch rest := res.(type) {
		case int64:
			return floatsEqual(float64(rest), v) || float64(rest) > v, nil
		case float64:
			return floatsEqual(rest, v) || rest > v, nil
		}
	}
	return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", op, fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
}

func compareString(res interface{}, op string, v string, opPos parsley.Pos) (interface{}, parsley.Error) {
	switch rest := res.(type) {
	case string:
		switch op {
		case "==":
			return rest == v, nil
		case "!=":
			return rest != v, nil
		case "<":
			return rest < v, nil
		case "<=":
			return rest <= v, nil
		case ">":
			return rest > v, nil
		case ">=":
			return rest >= v, nil
		}
	}
	return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", op, fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
}

func compareBool(res interface{}, op string, v bool, opPos parsley.Pos) (interface{}, parsley.Error) {
	switch rest := res.(type) {
	case bool:
		switch op {
		case "==":
			return rest == v, nil
		case "!=":
			return rest != v, nil
		}
	}
	return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", op, fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
}

func compareTime(res interface{}, op string, v time.Time, opPos parsley.Pos) (interface{}, parsley.Error) {
	switch rest := res.(type) {
	case time.Time:
		switch op {
		case "==":
			return rest == v, nil
		case "!=":
			return rest != v, nil
		case "<":
			return rest.Before(v), nil
		case "<=":
			return rest == v || rest.Before(v), nil
		case ">":
			return rest.After(v), nil
		case ">=":
			return rest == v || rest.After(v), nil
		}
	}
	return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", op, fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
}

func compareTimeDuration(res interface{}, op string, v time.Duration, opPos parsley.Pos) (interface{}, parsley.Error) {
	switch rest := res.(type) {
	case time.Duration:
		switch op {
		case "==":
			return rest == v, nil
		case "!=":
			return rest != v, nil
		case "<":
			return rest < v, nil
		case "<=":
			return rest <= v, nil
		case ">":
			return rest > v, nil
		case ">=":
			return rest >= v, nil
		}
	}
	return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", op, fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
}

func compareOther(res interface{}, op string, v interface{}, opPos parsley.Pos) (interface{}, parsley.Error) {
	if reflect.TypeOf(res) == reflect.TypeOf(v) {
		switch op {
		case "==":
			return reflect.DeepEqual(res, v), nil
		case "!=":
			return !reflect.DeepEqual(res, v), nil
		}
	}
	return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", op, fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
}
