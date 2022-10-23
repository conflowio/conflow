// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import (
	"fmt"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

var _ conflow.FunctionNode = &Node{}

// Node is a function node
type Node struct {
	nameNode      *conflow.NameNode
	argumentNodes []parsley.Node
	readerPos     parsley.Pos
	interpreter   conflow.FunctionInterpreter
	schema        schema.Schema
}

// Name returns with the function name
func (n *Node) Name() conflow.ID {
	return n.nameNode.Value().(conflow.ID)
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
	parameters := n.interpreter.Schema().(*schema.Function).GetParameters()

	args := make([]interface{}, 0, len(n.argumentNodes))

	for i := 0; i < len(parameters); i++ {
		v, evalErr := parsley.EvaluateNode(ctx, n.argumentNodes[i])
		if evalErr != nil {
			return nil, evalErr
		}
		vv, verr := parameters[i].Schema.ValidateValue(v)
		if verr != nil {
			return nil, parsley.NewError(n.argumentNodes[i].Pos(), verr)
		}
		args = append(args, vv)
	}

	if n.interpreter.Schema().(*schema.Function).AdditionalParameters != nil {
		for i := len(parameters); i < len(n.argumentNodes); i++ {
			v, evalErr := parsley.EvaluateNode(ctx, n.argumentNodes[i])
			if evalErr != nil {
				return nil, evalErr
			}
			vv, verr := n.interpreter.Schema().(*schema.Function).GetAdditionalParameters().Schema.ValidateValue(v)
			if verr != nil {
				return nil, parsley.NewError(n.argumentNodes[i].Pos(), verr)
			}
			args = append(args, vv)
		}
	}

	res, resErr := n.interpreter.Eval(ctx, args)
	if resErr != nil {
		if funcErr, ok := resErr.(*Error); ok {
			return nil, parsley.NewError(n.argumentNodes[funcErr.ArgIndex].Pos(), funcErr.Err)
		}
		return nil, parsley.NewError(n.Pos(), resErr)
	}
	return res, nil
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
