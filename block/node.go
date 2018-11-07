// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"errors"
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/identifier"
	"github.com/opsidian/parsley/parsley"
)

// Node is a block node
type Node struct {
	typeNode    *identifier.Node
	idNode      *identifier.Node
	paramNodes  map[basil.ID]basil.BlockParamNode
	blockNodes  []basil.BlockNode
	readerPos   parsley.Pos
	interpreter Interpreter
	blockType   string
}

// ID returns with the id of the block
func (n *Node) ID() basil.ID {
	if n.idNode == nil {
		return ""
	}

	id, _ := n.idNode.Value(nil)
	return id.(basil.ID)
}

// Token returns with the node's token
func (n *Node) Token() string {
	return "BLOCK"
}

// Type returns with the node's type
func (n *Node) Type() string {
	nodeType, _ := n.typeNode.Value(nil)
	return string(nodeType.(basil.ID))
}

// StaticCheck runs static analysis on the node
func (n *Node) StaticCheck(ctx interface{}) parsley.Error {
	if !n.interpreter.HasForeignID() {
		if n.idNode != nil {
			idRegistry := ctx.(basil.IDRegistryAware).GetIDRegistry()
			if err := idRegistry.RegisterID(n.ID()); err != nil {
				return parsley.NewError(n.idNode.Pos(), err)
			}
		}
	} else {
		if n.idNode == nil {
			return parsley.NewError(n.idNode.Pos(), errors.New("identifier must be set"))
		}
	}
	uniqueBlockIDs := map[string]struct{}{}
	for _, blockNode := range n.blockNodes {
		if _, exists := n.interpreter.NodeTransformer(blockNode.Type()); !exists {
			return parsley.NewErrorf(n.idNode.Pos(), "%q block type is invalid", blockNode.Type())
		}
		if blockNode.ID() != "" {
			blockID := fmt.Sprintf("%s.%s", blockNode.Type(), blockNode.ID())
			if _, exists := uniqueBlockIDs[blockID]; exists {
				return parsley.NewErrorf(blockNode.Pos(), "%q was defined multiple times", blockID)
			}
			uniqueBlockIDs[blockID] = struct{}{}
		}
	}

	_, err := n.interpreter.StaticCheck(ctx, n)
	if err != nil {
		return err
	}

	return nil
}

// Value creates a new block
func (n *Node) Value(ctx interface{}) (interface{}, parsley.Error) {
	if n.idNode == nil {
		idRegistry := ctx.(basil.IDRegistryAware).GetIDRegistry()
		id := idRegistry.GenerateID()
		n.idNode = identifier.NewNode(id, n.typeNode.ReaderPos(), n.typeNode.ReaderPos())
	}

	return n.interpreter.Eval(ctx, n)
}

// Eval evaluates the given stage on an existing block
func (n *Node) Eval(ctx interface{}, stage string, block basil.Block) parsley.Error {
	return n.interpreter.EvalBlock(ctx, n, stage, block)
}

// Pos returns with the node's position
func (n *Node) Pos() parsley.Pos {
	return n.typeNode.Pos()
}

// ReaderPos returns with the reader's position
func (n *Node) ReaderPos() parsley.Pos {
	return n.readerPos
}

// SetReaderPos amends the reader position using the given function
func (n *Node) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	n.readerPos = f(n.readerPos)
}

// ParamNodes returns with the parameter nodes
func (n *Node) ParamNodes() map[basil.ID]basil.BlockParamNode {
	return n.paramNodes
}

// BlockNodes returns with the child block nodes
func (n *Node) BlockNodes() []basil.BlockNode {
	return n.blockNodes
}
