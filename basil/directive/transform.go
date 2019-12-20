package directive

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

func Transform(
	parseCtx interface{},
	nodes []parsley.Node,
) ([]basil.BlockNode, basil.Dependencies, parsley.Error) {
	directives := make([]basil.BlockNode, 0, len(nodes))
	dependencies := make(basil.Dependencies)

	for _, n := range nodes {
		res, err := n.(parsley.Transformable).Transform(parseCtx)
		if err != nil {
			return nil, nil, err
		}
		blockNode := res.(basil.BlockNode)
		directives = append(directives, blockNode)

		parsley.Walk(blockNode, func(node parsley.Node) bool {
			if v, ok := node.(basil.VariableNode); ok {
				dependencies[v.ID()] = v
			}
			return false
		})
	}
	return directives, dependencies, nil
}
