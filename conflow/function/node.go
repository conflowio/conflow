// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

var _ conflow.FunctionNode = &Node{}

// Node is a function node
type Node struct {
	nameNode      parsley.Node
	argumentNodes []parsley.Node
	readerPos     parsley.Pos
	interpreter   conflow.FunctionInterpreter
	schema        schema.Schema
}

// Name returns with the function name
func (n *Node) Name() conflow.ID {
	value, _ := parsley.EvaluateNode(nil, n.nameNode)
	switch v := value.(type) {
	case conflow.ID:
		return v
	case []conflow.ID:
		return conflow.ID(fmt.Sprintf("%s.%s", v[0], v[1]))
	default:
		panic(fmt.Errorf("unexpected name node value: %T", v))
	}
}

// Token returns with the node's token
func (n *Node) Token() string {
	return "FUNC"
}

// Schema returns the schema for the node's value
func (n *Node) Schema() interface{} {
	return n.schema
}

// StaticCheck runs static analysis on the node
func (n *Node) StaticCheck(ctx interface{}) parsley.Error {
	s := n.interpreter.Schema().(*schema.Function)
	name := n.Name()

	if len(n.argumentNodes) < len(s.Parameters) {
		pos := n.nameNode.Pos()
		if len(n.argumentNodes) > 0 {
			pos = n.argumentNodes[len(n.argumentNodes)-1].ReaderPos()
		}
		countReq := "exactly"
		if s.AdditionalParameters != nil {
			countReq = "at least"
		}
		argumentStr := "argument"
		if len(s.Parameters) > 1 {
			argumentStr = "arguments"
		}
		return parsley.NewErrorf(
			pos,
			"%s requires %s %d %s, but got %d",
			name,
			countReq,
			len(s.Parameters),
			argumentStr,
			len(n.argumentNodes),
		)
	}

	if len(n.argumentNodes) > len(s.Parameters) && s.AdditionalParameters == nil {
		argumentStr := "argument"
		if len(s.Parameters) > 1 {
			argumentStr = "arguments"
		}
		return parsley.NewErrorf(
			n.argumentNodes[len(s.Parameters)].Pos(),
			"%s requires exactly %d %s, but got %d",
			name,
			len(s.Parameters),
			argumentStr,
			len(n.argumentNodes),
		)
	}

	n.schema = s.Result
	if n.schema == nil {
		panic("a function must have a return value")
	}

	for i, arg := range n.argumentNodes {
		var paramSchema schema.Schema
		if i < len(s.Parameters) {
			paramSchema = s.Parameters[i].Schema
		} else {
			paramSchema = s.AdditionalParameters.Schema
		}

		if err := paramSchema.ValidateSchema(arg.Schema().(schema.Schema), false); err != nil {
			return parsley.NewError(arg.Pos(), err)
		}

		if s.ResultTypeFrom != "" && i < len(s.Parameters) {
			if s.ResultTypeFrom == s.Parameters[i].Name {
				n.schema = arg.Schema().(schema.Schema)
			}
		}
	}

	return nil
}

// Value returns with the result of the function
func (n *Node) Value(ctx interface{}) (interface{}, parsley.Error) {
	return n.interpreter.Eval(ctx, n)
}

// Pos returns with the node's position
func (n *Node) Pos() parsley.Pos {
	return n.nameNode.Pos()
}

// ReaderPos returns with the reader's position
func (n *Node) ReaderPos() parsley.Pos {
	return n.readerPos
}

// SetReaderPos amends the reader position using the given function
func (n *Node) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	n.readerPos = f(n.readerPos)
}

// ArgumentNodes returns with the function argument nodes
func (n *Node) ArgumentNodes() []parsley.Node {
	return n.argumentNodes
}

// Children returns with the argument nodes
func (n *Node) Children() []parsley.Node {
	return n.argumentNodes
}

func (n *Node) String() string {
	return fmt.Sprintf("%s{%s, %s, %d..%d}", n.Token(), n.Name(), n.argumentNodes, n.Pos(), n.ReaderPos())
}
