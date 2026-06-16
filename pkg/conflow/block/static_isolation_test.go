// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"context"
	"testing"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/conflowfakes"
	"github.com/conflowio/conflow/pkg/conflow/parameter"
	"github.com/conflowio/conflow/pkg/loggers/zerolog"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/values"
)

func TestStaticEvaluateChildBindParameter(t *testing.T) {
	var captured interface{}

	arraySchema := &schema.Array{Items: &schema.String{}}
	objectSchema := &schema.Object{
		ParameterNames: map[string]string{"Items": "items"},
		Properties: map[string]schema.Schema{
			"Items": arraySchema,
		},
	}

	upstream := []interface{}{"shared"}
	valueNode := &conflowfakes.FakeNode{}
	valueNode.ValueReturns(upstream, nil)

	interpreter := &conflowfakes.FakeBlockInterpreter{}
	interpreter.SchemaReturns(objectSchema)
	blockInstance := &conflowfakes.FakeIdentifiable{}
	interpreter.CreateBlockReturns(blockInstance)
	interpreter.SetParamStub = func(_ conflow.Block, _ conflow.ID, value interface{}) error {
		captured = value
		return nil
	}

	blockNode := &conflowfakes.FakeBlockNode{}
	blockNode.IDReturns("parent")
	blockNode.InterpreterReturns(interpreter)

	logger := zerolog.NewDisabledLogger()
	evalCtx := conflow.NewEvalContext(context.Background(), nil, logger, nil, nil, blockNode)

	staticContainer := NewStaticContainer(evalCtx, blockNode)
	staticContainer.block = blockInstance

	paramNode := parameter.NewNode(
		"parent",
		conflow.NewIDNode("items", conflow.ClassifierNone, parsley.NilPos, parsley.NilPos),
		valueNode,
		true,
		nil,
	)
	paramNode.SetSchema(arraySchema)

	if err := staticContainer.evaluateChild(paramNode); err != nil {
		t.Fatalf("evaluateChild: %v", err)
	}

	list, ok := captured.(*values.List[interface{}])
	if !ok {
		t.Fatalf("SetParam value type = %T, want *values.List[interface{}]", captured)
	}
	if list.At(0) != "shared" {
		t.Fatalf("bound list[0] = %q, want %q", list.At(0), "shared")
	}

	upstream[0] = "mutated"
	if list.At(0) != "shared" {
		t.Fatalf("bound list mutated after upstream change")
	}
}
