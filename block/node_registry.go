package block

import (
	"fmt"

	"github.com/opsidian/basil/basil"
)

// NodeRegistry stores block nodes
type NodeRegistry map[basil.ID]basil.BlockNode

// NewNodeRegistry creates a new node registry
func NewNodeRegistry() NodeRegistry {
	return NodeRegistry(make(map[basil.ID]basil.BlockNode, 0))
}

// BlockNode returns with the given block node if it exists
func (n NodeRegistry) BlockNode(id basil.ID) (basil.BlockNode, bool) {
	node, ok := n[id]
	return node, ok
}

// AddBlockNode adds a new block node
// It returns with an error if a block with the same id was already registered
func (n NodeRegistry) AddBlockNode(node basil.BlockNode) error {
	if _, exists := n[node.ID()]; exists {
		return fmt.Errorf("duplicated identifier: %q", node.ID())
	}

	n[node.ID()] = node

	return nil
}
