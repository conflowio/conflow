package basil

import "github.com/opsidian/parsley/parsley"

// VariableNode stores a variable reference. It always refers to a named block's parameter.
//go:generate counterfeiter . VariableNode
type VariableNode interface {
	parsley.Node
	Identifiable
	ParentID() ID
	ParamName() ID
}
