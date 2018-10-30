package ocl

import (
	"github.com/opsidian/parsley/parsley"
)

// Block is an interface for a block object
type Block interface {
	Identifiable
	TypeAware
	ContextAware
}

// BlockPreInterpreter as an interface for running a pre-evaluation hook
type BlockPreInterpreter interface {
	PreEval(ctx interface{}, stage string) parsley.Error
}

// BlockPostInterpreter as an interface for running a post-evaluation hook
type BlockPostInterpreter interface {
	PostEval(ctx interface{}, stage string) parsley.Error
}
