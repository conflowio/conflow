// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"fmt"
)

// Evaluate will evaluate the given node which was previously parsed with the passed parse context
func Evaluate(
	parseCtx *ParseContext,
	blockContext BlockContext,
	scheduler Scheduler,
	id ID,
) (interface{}, error) {
	node, ok := parseCtx.BlockNode(id)
	if !ok {
		return nil, fmt.Errorf("block %q does not exist", id)
	}

	value, err := node.Value(NewEvalContext(blockContext, scheduler))
	if err != nil {
		return nil, parseCtx.FileSet().ErrorWithPosition(err)
	}

	return value, nil
}
