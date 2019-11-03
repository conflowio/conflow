// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Element will match a variable expression defined by the following rule, where P is the input parser:
//   S         -> P (VAR_INDEX)*
//   VAR_INDEX -> "." ID
//             -> "[" P "]"
//   ID        -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
func Element(p parsley.Parser, index parsley.Parser) parser.Func {
	arrayIndex := combinator.SeqOf(
		terminal.Rune('['),
		text.LeftTrim(index, text.WsSpaces),
		text.LeftTrim(terminal.Rune(']'), text.WsSpaces),
	).Bind(interpreter.Select(1))

	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return p
		}
		return arrayIndex
	}
	lenCheck := func(len int) bool {
		return len > 0
	}
	return combinator.Single(
		combinator.Seq("VAR", lookup, lenCheck).Bind(ast.InterpreterFunc(evalElement)),
	)
}

func evalElement(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	res, err := nodes[0].Value(ctx)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nodes); i++ {
		index, err := nodes[i].Value(ctx)
		if err != nil {
			return nil, err
		}
		switch rest := res.(type) {
		case []interface{}:
			switch indext := index.(type) {
			case int64:
				if indext >= 0 && indext < int64(len(rest)) {
					res = rest[indext]
				} else {
					indexNode := nodes[i].(parsley.NonTerminalNode).Children()[1]
					return nil, parsley.NewErrorf(indexNode.Pos(), "array index out of bounds: %d (0..%d)", indext, len(rest)-1)
				}
			default:
				indexNode := nodes[i].(parsley.NonTerminalNode).Children()[1]
				return nil, parsley.NewErrorf(indexNode.Pos(), "non-integer index on array")
			}
		case map[string]interface{}:
			switch indext := index.(type) {
			case string:
				val, ok := rest[indext]
				if !ok {
					indexNode := nodes[i].(parsley.NonTerminalNode).Children()[1]
					return nil, parsley.NewErrorf(indexNode.Pos(), "key %q does not exist on map", indext)
				}
				res = val
			default:
				indexNode := nodes[i].(parsley.NonTerminalNode).Children()[1]
				return nil, parsley.NewErrorf(indexNode.Pos(), "invalid non-string index on map")
			}
		default:
			indexNode := nodes[i].(parsley.NonTerminalNode).Children()[1]
			return nil, parsley.NewErrorf(indexNode.Pos(), "can not get index on %s type", fmt.Sprintf("%T", res))
		}
	}
	return res, nil
}
