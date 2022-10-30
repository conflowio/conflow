// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"strings"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/util/validation"
)

// TransformPathErrors will try to transform errors which a path
// The exact node position will be looked up by the path expression separated by "." characters, where a part element can be
//   - parameter
//   - parameter[1]
//   - parameter["key"]
func TransformPathErrors(
	parseCtx *ParseContext,
	err parsley.Error,
	node parsley.Node,
) error {
	te, ok := err.Cause().(interface {
		TransformError(func(path string, err error) error) error
	})
	if !ok {
		return parseCtx.FileSet().ErrorWithPosition(err)
	}

	return te.TransformError(func(path string, err error) error {
		res := func() parsley.Error {
			if pe, ok := err.(parsley.Error); ok {
				return pe
			}

			if path == "" {
				return parsley.NewError(node.Pos(), err)
			}

			n, remainingPath := FindNodeByPath(node, strings.Split(path, "."))
			if len(remainingPath) == 0 {
				return parsley.NewError(n.Pos(), err)
			} else {
				return parsley.NewError(n.Pos(), validation.NewFieldError(strings.Join(remainingPath, "."), err))
			}
		}()
		return parseCtx.FileSet().ErrorWithPosition(res)
	})
}
