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

	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/schema"
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

	o := node.Interpreter().Schema().(*schema.Object)
	for jsonPropertyName, property := range o.Properties {
		if property.GetAnnotation(annotations.UserDefined) == "true" && o.IsPropertyRequired(jsonPropertyName) {
			parameterName := o.ParameterName(jsonPropertyName)
			if _, isDefined := inputParams[ID(parameterName)]; !isDefined {
				return nil, fmt.Errorf("%q input parameter must be defined", parameterName)
			}
		}
	}

	for k, v := range inputParams {
		property, ok := o.PropertyByParameterName(string(k))
		if ok && property.GetAnnotation(annotations.UserDefined) == "true" && !property.GetReadOnly() {
			nv, err := property.ValidateValue(v)
			if err != nil {
				return nil, fmt.Errorf("invalid input parameter %q: %w", k, err)
			}
			inputParams[k] = nv
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
