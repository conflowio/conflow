// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Variable will match a variable expression defined by the following rule, where P is the input parser:
//   S         -> ID (VAR_INDEX)*
//   ID        -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   VAR_INDEX -> "." ID
//             -> "[" P "]"
func Variable(p parsley.Parser) *combinator.Sequence {
	variableIndex := combinator.SeqOf(
		terminal.Rune('.'),
		ID(),
	).Bind(interpreter.Select(1))
	arrayIndex := combinator.SeqOf(
		terminal.Rune('['),
		text.LeftTrim(p, text.WsSpaces),
		text.LeftTrim(terminal.Rune(']'), text.WsSpaces),
	).Bind(interpreter.Select(1))

	id := ID()
	arrayOrVariableIndex := combinator.Choice(arrayIndex, variableIndex)

	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return id
		}
		return arrayOrVariableIndex
	}
	lenCheck := func(len int) bool {
		return len > 0
	}
	return combinator.Seq("VAR", lookup, lenCheck).Bind(ast.InterpreterFunc(evalVariable))
}

func evalVariable(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
	value0, _ := nodes[0].Value(ctx)
	name := value0.(string)
	varIndex := make([]interface{}, 0, len(nodes)-1)
	for i := 1; i < len(nodes); i++ {
		val, err := nodes[i].Value(ctx)
		if err != nil {
			return nil, err
		}
		varIndex = append(varIndex, val)
	}

	variableProvider := ctx.(basil.VariableProviderAware).GetVariableProvider()

	res, err := variableProvider.LookupVar(lookup(name, varIndex, nodes))
	if err != nil {
		varName := name
		for _, index := range varIndex {
			varName = varName + "[" + fmt.Sprintf("%v", index) + "]"
		}
		if err == basil.ErrVariableNotFound {
			return nil, parsley.WrapError(
				parsley.NewError(nodes[0].Pos(), err),
				"variable '%s' does not exist", varName,
			)
		}

		return nil, parsley.NewError(nodes[0].Pos(), err)
	}
	return res, nil
}

func lookup(name string, varIndex []interface{}, nodes []parsley.Node) basil.VariableLookUp {
	return func(provider basil.VariableProvider) (interface{}, error) {
		res, ok := provider.GetVar(name)
		if !ok {
			return nil, basil.ErrVariableNotFound
		}
		for i, index := range varIndex {
			switch rest := res.(type) {
			case []interface{}:
				switch indext := index.(type) {
				case int64:
					if indext >= 0 && indext < int64(len(rest)) {
						res = rest[indext]
					} else {
						indexNode := nodes[i+1].(parsley.NonTerminalNode).Children()[1]
						return nil, parsley.NewErrorf(indexNode.Pos(), "array index out of bounds: %d (0..%d)", indext, len(rest)-1)
					}
				default:
					indexNode := nodes[i+1].(parsley.NonTerminalNode).Children()[1]
					return nil, parsley.NewErrorf(indexNode.Pos(), "invalid non-integer index on array")
				}
			case map[string]interface{}:
				switch indext := index.(type) {
				case string:
					res, ok = rest[indext]
					if !ok {
						return nil, basil.ErrVariableNotFound
					}
				default:
					indexNode := nodes[i+1].(parsley.NonTerminalNode).Children()[1]
					return nil, parsley.NewErrorf(indexNode.Pos(), "invalid non-string index on map")
				}
			default:
				indexNode := nodes[i+1].(parsley.NonTerminalNode).Children()[1]
				return nil, parsley.NewErrorf(indexNode.Pos(), "can not get index on %s type", fmt.Sprintf("%T", res))
			}
		}
		return res, nil
	}
}
