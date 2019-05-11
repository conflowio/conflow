package block

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// InterpreterRegistry contains a list of block interpreters and behaves as a node transformer registry
type InterpreterRegistry map[string]basil.BlockInterpreter

// NodeTransformer returns with the named node transformer
func (i InterpreterRegistry) NodeTransformer(name string) (parsley.NodeTransformer, bool) {
	interpreter, exists := i[name]
	if !exists {
		return nil, false
	}

	if name == basil.MainID {
		return parsley.NodeTransformFunc(func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
			return transformMainNode(userCtx, node, interpreter)
		}), true
	}

	return parsley.NodeTransformFunc(func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
		return transformNode(userCtx, node, interpreter)
	}), true
}
