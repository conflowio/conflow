// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"context"
	"fmt"
	"reflect"
)

// Evaluate will evaluate the given node which was previously parsed with the passed parse context
func Evaluate(
	parseCtx *ParseContext,
	context context.Context,
	userContext interface{},
	logger Logger,
	scheduler JobScheduler,
	id ID,
	inputParams map[ID]interface{},
) (interface{}, error) {
	node, ok := parseCtx.BlockNode(id)
	if !ok {
		return nil, fmt.Errorf("block %q does not exist", id)
	}

	for paramName, param := range node.Interpreter().Params() {
		if param.IsUserDefined && param.IsRequired {
			if _, isDefined := inputParams[paramName]; !isDefined {
				return nil, fmt.Errorf("%q input parameter must be defined", paramName)
			}
		}
	}

	for k, v := range inputParams {
		param := node.Interpreter().Params()[k]
		if param.IsUserDefined && !param.IsOutput {
			if reflect.TypeOf(v).Name() != param.Type {
				return nil, fmt.Errorf("invalid input parameter %q, type must be %s", k, param.Type)
			}
		} else {
			return nil, fmt.Errorf("unknown input parameter: %q", k)
		}
	}

	evalContext := NewEvalContext(context, userContext, logger, scheduler, nil)
	evalContext.InputParams = inputParams

	value, err := node.Value(evalContext)
	if err != nil {
		return nil, parseCtx.FileSet().ErrorWithPosition(err)
	}

	return value, nil
}
