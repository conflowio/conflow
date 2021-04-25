// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// Node tokens
const (
	TokenBlock     = "BLOCK"
	TokenDirective = "BLOCK_DIRECTIVE"
	TokenBlockBody = "BLOCK_BODY"
)

var _ basil.BlockNode = &Node{}

// Node is a block node
type Node struct {
	typeNode     *basil.IDNode
	idNode       *basil.IDNode
	children     []basil.Node
	token        string
	directives   []basil.BlockNode
	readerPos    parsley.Pos
	interpreter  basil.BlockInterpreter
	dependencies basil.Dependencies
	provides     []basil.ID
	generates    []basil.ID
	schema       *schema.Object
}

// NewNode creates a new block node
func NewNode(
	idNode *basil.IDNode,
	typeNode *basil.IDNode,
	children []basil.Node,
	token string,
	directives []basil.BlockNode,
	readerPos parsley.Pos,
	interpreter basil.BlockInterpreter,
	dependencies basil.Dependencies,
) *Node {
	var provides []basil.ID
	var generates []basil.ID
	for _, c := range children {
		if b, ok := c.(basil.BlockNode); ok {
			if b.Schema().(schema.Schema).GetAnnotation("generated") == "true" {
				generates = append(generates, b.ID())
				generates = append(generates, b.Provides()...)
			} else {
				provides = append(provides, b.ID())
				provides = append(provides, b.Provides()...)
			}
			generates = append(generates, b.Generates()...)
		}
	}

	return &Node{
		idNode:       idNode,
		typeNode:     typeNode,
		children:     children,
		token:        token,
		directives:   directives,
		interpreter:  interpreter,
		readerPos:    readerPos,
		dependencies: dependencies,
		generates:    generates,
		provides:     provides,
		schema:       interpreter.Schema().(*schema.Object),
	}
}

// ID returns with the id of the block
func (n *Node) ID() basil.ID {
	return n.idNode.ID()
}

// BlockType returns with the type of block
func (n *Node) BlockType() basil.ID {
	return n.typeNode.ID()
}

// Token returns with the node's token
func (n *Node) Token() string {
	return n.token
}

// Schema returns the schema for the node's value
func (n *Node) Schema() interface{} {
	return n.schema
}

func (n *Node) SetSchema(s schema.Schema) {
	n.schema = n.interpreter.Schema().Copy().(*schema.Object)
	switch st := s.(type) {
	case *schema.Array:
		n.schema.Metadata.Merge(st.Items.(*schema.Reference).Metadata)
	case *schema.Reference:
		n.schema.Metadata.Merge(st.Metadata)
	default:
		panic(fmt.Errorf("unexpected schema for a block node: %T", s))
	}
}

func (n *Node) GetPropertySchema(name basil.ID) (schema.Schema, bool) {
	s, ok := n.schema.Properties[string(name)]
	if ok {
		return s, true
	}

	for _, n := range n.children {
		if p, ok := n.(basil.ParameterNode); ok && p.IsDeclaration() {
			if p.Name() == name {
				return p.Schema().(schema.Schema), true
			}
		}
	}

	return nil, false
}

// EvalStage returns with the evaluation stage
func (n *Node) EvalStage() basil.EvalStage {
	evalStageStr := n.schema.GetAnnotation("eval_stage")
	if evalStageStr != "" {
		return basil.EvalStages[evalStageStr]
	}

	return basil.EvalStageMain
}

// Dependencies returns the blocks/parameters this block depends on
func (n *Node) Dependencies() basil.Dependencies {
	return n.dependencies
}

// Interpreter returns with the interpreter
func (n *Node) Interpreter() basil.BlockInterpreter {
	return n.interpreter
}

// Generates returns true if any of the child blocks are generated
func (n *Node) Generates() []basil.ID {
	return n.generates
}

// Provides returns with the all the defined blocked node ids inside this block
func (n *Node) Provides() []basil.ID {
	return n.provides
}

// StaticCheck runs static analysis on the node
func (n *Node) StaticCheck(ctx interface{}) parsley.Error {
	blockCounts := map[basil.ID]int{}

	for _, child := range n.Children() {
		switch c := child.(type) {
		case basil.BlockNode:
			blockCounts[c.BlockType()] = blockCounts[c.BlockType()] + 1
			if property, ok := n.schema.Properties[string(c.BlockType())]; ok {
				if blockCounts[c.BlockType()] > 1 && property.Type() != schema.TypeArray {
					return parsley.NewError(c.Pos(), fmt.Errorf("%q block can only be defined once", c.BlockType()))
				}
			}
		case basil.ParameterNode:
			property, exists := n.schema.Properties[string(c.Name())]

			switch {
			case exists && c.IsDeclaration() && !schema.HasAnnotationValue(property, "user_defined", "true"):
				return parsley.NewErrorf(c.Pos(), "%q parameter already exists. Use \"=\" to set the parameter value or use a different name", c.Name())
			case exists && !c.IsDeclaration() && schema.HasAnnotationValue(property, "user_defined", "true"):
				return parsley.NewErrorf(c.Pos(), "%q must be defined as a new variable using \":=\"", c.Name())
			case !exists && !c.IsDeclaration():
				return parsley.NewErrorf(c.Pos(), "%q parameter does not exist", c.Name())
			case !c.IsDeclaration() && !schema.HasAnnotationValue(property, "user_defined", "true") && property.GetReadOnly():
				return parsley.NewErrorf(c.Pos(), "%q is an read-only parameter and can not be set", c.Name())
			}
		default:
			panic(fmt.Errorf("invalid node type: %T", child))
		}
	}

	for _, required := range n.schema.Required {
		found := func() bool {
			for _, child := range n.Children() {
				switch c := child.(type) {
				case basil.BlockNode:
					if string(c.BlockType()) == required {
						return true
					}
				case basil.ParameterNode:
					if string(c.Name()) == required && c.ValueNode().Schema().(schema.Schema).Type() != schema.TypeNull {
						return true
					}
				}
			}
			return false
		}()
		if !found {
			if IsBlockSchema(n.schema.Properties[required]) {
				return parsley.NewError(n.Pos(), fmt.Errorf("%q block is required", required))
			} else {
				return parsley.NewError(n.Pos(), fmt.Errorf("%q parameter is required", required))
			}
		}
	}

	return nil
}

// Value creates a new block
func (n *Node) Value(userCtx interface{}) (interface{}, parsley.Error) {
	var container basil.JobContainer
	switch {
	case n.EvalStage() == basil.EvalStageParse || n.Token() == TokenDirective:
		container = NewStaticContainer(userCtx.(*basil.EvalContext), n)
	case n.Token() == TokenBlock || n.Token() == TokenBlockBody:
		container = NewContainer(userCtx.(*basil.EvalContext), n, nil, nil, nil, false)
	default:
		panic(fmt.Errorf("unknown block type: %s", n.Token()))
	}

	container.Run()
	return container.Value()
}

// Pos returns with the node's position
func (n *Node) Pos() parsley.Pos {
	return n.idNode.Pos()
}

// ReaderPos returns with the reader's position
func (n *Node) ReaderPos() parsley.Pos {
	return n.readerPos
}

// SetReaderPos amends the reader position using the given function
func (n *Node) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	n.readerPos = f(n.readerPos)
}

// Children returns with the parameter and child block nodes
func (n *Node) Children() []basil.Node {
	return n.children
}

// Directives returns with the directive blocks
func (n *Node) Directives() []basil.BlockNode {
	return n.directives
}

// Walk runs the given function on all child nodes
func (n *Node) Walk(f func(n parsley.Node) bool) bool {
	for _, node := range n.children {
		if parsley.Walk(node, f) {
			return true
		}
	}

	return false
}

func (n *Node) CreateContainer(
	ctx *basil.EvalContext,
	parent basil.BlockContainer,
	value interface{},
	wgs []basil.WaitGroup,
	pending bool,
) basil.JobContainer {
	return NewContainer(ctx, n, parent, value, wgs, pending)
}

func (n *Node) String() string {
	return fmt.Sprintf("%s{%s, %s, %s, %d..%d}", n.Token(), n.typeNode, n.idNode, n.children, n.Pos(), n.ReaderPos())
}
