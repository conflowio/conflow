// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"context"
	"fmt"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/conflow/schema"
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

	o := node.Interpreter().Schema().(schema.ObjectKind)
	for paramName, param := range o.GetProperties() {
		if schema.HasAnnotationValue(param, "user_defined", "true") && o.IsPropertyRequired(paramName) {
			if _, isDefined := inputParams[ID(paramName)]; !isDefined {
				return nil, fmt.Errorf("%q input parameter must be defined", paramName)
			}
		}
	}

	for k, v := range inputParams {
		property := o.GetProperties()[string(k)]
		if property != nil &&
			schema.HasAnnotationValue(property, "user_defined", "true") &&
			!property.GetReadOnly() {
			if err := property.ValidateValue(v); err != nil {
				return nil, fmt.Errorf("invalid input parameter %q: %w", k, err)
			}
		} else {
			return nil, fmt.Errorf("unknown input parameter: %q", k)
		}
	}

	evalContext := NewEvalContext(context, userContext, logger, scheduler, nil)
	evalContext.InputParams = inputParams

	value, err := parsley.EvaluateNode(evalContext, node)
	if err != nil {
		return nil, parseCtx.FileSet().ErrorWithPosition(err)
	}

	return value, nil
}
