package parser

import (
	"errors"
	"fmt"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// IDRegExp is the regular expression for a valid identifier
const IDRegExp = "[a-z][a-z0-9]*(?:_[a-z0-9]+)*"

// ID parses an identifier:
//   S -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//
// An ID can only contain lowercase letters, numbers and underscore characters.
// It must start with a letter and no duplicate underscores are allowed.
//
func ID() parser.Func {
	notFoundErr := errors.New("was expecting identifier")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, match := tr.ReadRegexp(pos, IDRegExp); match != nil {
			return NewIDNode(string(match), pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}

// IDNode represents an identifier
// If it is initialised with an empty string it will generate a value using a block registry
type IDNode struct {
	value     string
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewIDNode creates a new ID node
func NewIDNode(value string, pos parsley.Pos, readerPos parsley.Pos) *IDNode {
	return &IDNode{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (i *IDNode) Token() string {
	return "ID"
}

// Value returns with the value of the node
func (i *IDNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	if i.value == "" {
		blockRegistry := ctx.(ocl.BlockRegistryAware).GetBlockRegistry()
		i.value = blockRegistry.GenerateBlockID()
	}
	return i.value, nil
}

// Pos returns the position
func (i *IDNode) Pos() parsley.Pos {
	return i.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (i *IDNode) ReaderPos() parsley.Pos {
	return i.readerPos
}

// SetReaderPos changes the reader position
func (i *IDNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	i.readerPos = f(i.readerPos)
}

// String returns with a string representation of the node
func (i *IDNode) String() string {
	return fmt.Sprintf("ID{%v, %d..%d}", i.value, i.pos, i.readerPos)
}
