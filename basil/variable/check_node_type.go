// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package variable

import (
	"fmt"

	"github.com/opsidian/parsley/text/terminal"

	"github.com/opsidian/parsley/parsley"
)

// CheckNodeType checks the type of the node
// It will only return an error if it's sure the type is incorrent
func CheckNodeType(node parsley.Node, expectedType string) parsley.Error {
	if node.Type() == expectedType {
		return nil
	}

	if expectedType == TypeUnknown || expectedType == TypeAny {
		return nil
	}

	if node.Type() == TypeUnknown || node.Type() == TypeAny {
		return nil
	}

	if expectedType == TypeTime {
		if node.Type() == TypeString {
			if stringNode, ok := node.(*terminal.StringNode); ok {
				val, _ := stringNode.Value(nil)
				if _, err := TimeValue(val); err != nil {
					return parsley.NewError(node.Pos(), err)
				}
			}
			return nil
		}
	}

	if expectedType == TypeStringArray && node.Type() == TypeArray {
		for _, child := range node.(parsley.NonTerminalNode).Children() {
			if err := CheckNodeType(child, TypeString); err != nil {
				return err
			}
		}
		return nil
	}

	for unionType, types := range UnionTypes {
		if expectedType == unionType {
			for _, t := range types {
				if t == node.Type() {
					return nil
				}
			}
		}
	}

	typeErr, ok := TypeErrors[expectedType]
	if !ok {
		panic(fmt.Sprintf("Unknown type: %s", expectedType))
	}
	return parsley.NewError(node.Pos(), typeErr)
}
