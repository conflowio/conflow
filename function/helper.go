// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import (
	"reflect"
	"strings"

	"github.com/opsidian/parsley/parsley"
)

func evaluateParams(ctx interface{}, function parsley.Node, nodes []parsley.Node, kinds [][]reflect.Kind) ([]interface{}, parsley.Error) {
	err := checkParamCount(ctx, function, len(nodes), len(kinds))
	if err != nil {
		return nil, err
	}

	params := make([]interface{}, len(nodes))
	for i, node := range nodes {
		var evalErr parsley.Error
		if params[i], evalErr = evaluateParam(ctx, function, node, i, kinds[i]); evalErr != nil {
			return nil, evalErr
		}
	}
	return params, nil
}

func evaluateParam(ctx interface{}, function parsley.Node, param parsley.Node, index int, kinds []reflect.Kind) (interface{}, parsley.Error) {
	value, err := param.Value(ctx)
	if err != nil {
		return nil, err
	}

	valueKind := reflect.ValueOf(value).Kind()
	for _, kind := range kinds {
		if kind == valueKind {
			return value, nil
		}
	}
	functionName, _ := function.Value(ctx)

	if len(kinds) == 1 {
		return nil, parsley.NewErrorf(param.Pos(), "%s expects %s as %d. parameter, got %s", functionName, kinds[0], index+1, reflect.ValueOf(value).Type())
	}

	kindStr := make([]string, len(kinds))
	for i, kind := range kinds {
		kindStr[i] = kind.String()
	}
	return nil, parsley.NewErrorf(param.Pos(), "%s expects %s as %d. parameter, got %s", functionName, strings.Join(kindStr, "|"), index+1, reflect.ValueOf(value).Kind())
}

func checkParamCount(ctx interface{}, function parsley.Node, paramCount int, expectedCount int) parsley.Error {
	if paramCount != expectedCount {
		functionName, _ := function.Value(ctx)
		pStr := "parameters"
		if expectedCount == 1 {
			pStr = "parameter"
		}
		return parsley.NewErrorf(function.Pos(), "%s expects %d %s, got %d", functionName, expectedCount, pStr, paramCount)
	}
	return nil
}
