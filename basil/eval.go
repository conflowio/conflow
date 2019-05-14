package basil

import (
	"context"
	"fmt"
)

// Evaluate will evaluate the given node which was previously parsed with the passed parse context
func Evaluate(
	parseCtx *ParseContext,
	blockContext BlockContextOverride,
	id ID,
) (interface{}, error) {
	node, ok := parseCtx.BlockNode(id)
	if !ok {
		return nil, fmt.Errorf("block %q does not exist", id)
	}

	if blockContext.Context == nil {
		blockContext.Context = context.Background()
	}

	value, err := node.Value(NewEvalContext(blockContext.Context, blockContext.UserContext))
	if err != nil {
		return nil, parseCtx.FileSet().ErrorWithPosition(err)
	}

	return value, nil
}
