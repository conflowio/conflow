// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directive

import (
	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
)

func Transform(
	parseCtx interface{},
	nodes []parsley.Node,
) ([]conflow.BlockNode, conflow.Dependencies, parsley.Error) {
	directives := make([]conflow.BlockNode, 0, len(nodes))
	dependencies := make(conflow.Dependencies)

	for _, n := range nodes {
		res, err := n.(parsley.Transformable).Transform(parseCtx)
		if err != nil {
			return nil, nil, err
		}
		blockNode := res.(conflow.BlockNode)
		directives = append(directives, blockNode)

		parsley.Walk(blockNode, func(node parsley.Node) bool {
			if v, ok := node.(conflow.VariableNode); ok {
				dependencies[v.ID()] = v
			}
			return false
		})
	}
	return directives, dependencies, nil
}
