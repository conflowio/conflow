package identifier

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// Node represents an identifier
// If it is initialised with an empty string it will generate a value using a block registry
type Node struct {
	value     basil.ID
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewNode creates a new ID node
func NewNode(value basil.ID, pos parsley.Pos, readerPos parsley.Pos) *Node {
	return &Node{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (n *Node) Token() string {
	return "ID"
}

// Type returns with the type of the node
func (n *Node) Type() string {
	return "ID"
}

// Value returns with the value of the node
func (n *Node) Value(ctx interface{}) (interface{}, parsley.Error) {
	if n.value == "" {
		idRegistry := ctx.(basil.IDRegistryAware).GetIDRegistry()
		n.value = idRegistry.GenerateID()
	}
	return n.value, nil
}

// Pos returns the position
func (n *Node) Pos() parsley.Pos {
	return n.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (n *Node) ReaderPos() parsley.Pos {
	return n.readerPos
}

// SetReaderPos changes the reader position
func (n *Node) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	n.readerPos = f(n.readerPos)
}

// String returns with a string representation of the node
func (n *Node) String() string {
	return fmt.Sprintf("ID{%v, %d..%d}", n.value, n.pos, n.readerPos)
}
