package basil

import "github.com/opsidian/parsley/parsley"

// Node is an identifiable node which has dependencies and has an evaluation stage
//go:generate counterfeiter . Node
type Node interface {
	parsley.Node
	ID() ID
	EvalStage() EvalStage
	Dependencies() []IdentifiableNode
	Provides() []ID
}
