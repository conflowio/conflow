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

func TestSetChildBindSchemaParameter(t *testing.T) {
	fixture := newSetChildFixture(t)

	upstream := []interface{}{"shared"}
	if err := setChildWithValue(fixture.container, "items", upstream); err != nil {
		t.Fatalf("setChild: %v", err)
	}

	list, ok := (*fixture.setParamValue).(*values.List[interface{}])
	if !ok {
		t.Fatalf("SetParam value type = %T, want *values.List[interface{}]", *fixture.setParamValue)
	}
	if list.At(0) != "shared" {
		t.Fatalf("bound list[0] = %q, want %q", list.At(0), "shared")
	}

	upstream[0] = "mutated"
	if list.At(0) != "shared" {
		t.Fatalf("bound list mutated after upstream change")
	}
}

func TestSetChildBindExtraParams(t *testing.T) {
	fixture := newSetChildFixture(t)
	fixture.container.node.Interpreter().(*conflowfakes.FakeBlockInterpreter).SchemaReturns(&schema.Object{
		ParameterNames: map[string]string{},
		Properties:     map[string]schema.Schema{},
	})

	upstream := []interface{}{"extra"}
	if err := setChildWithValue(fixture.container, "extra_param", upstream); err != nil {
		t.Fatalf("setChild: %v", err)
	}
	if *fixture.setParamValue != nil {
		t.Fatal("expected extra param to bypass SetParam")
	}

	stored := fixture.container.extraParams["extra_param"]
	list, ok := stored.(*values.List[interface{}])
	if !ok {
		t.Fatalf("extra param type = %T, want *values.List[interface{}]", stored)
	}
	if list.At(0) != "extra" {
		t.Fatalf("extra param[0] = %q, want %q", list.At(0), "extra")
	}

	upstream[0] = "mutated"
	if list.At(0) != "extra" {
		t.Fatalf("extra param mutated after upstream change")
	}
}

func TestSetChildBindError(t *testing.T) {
	fixture := newSetChildFixture(t)
	if err := setChildWithValue(fixture.container, "items", "not-an-array"); err == nil {
		t.Fatal("expected bind error, got nil")
	}
}

type setChildFixture struct {
	setParamValue *interface{}
	container     *Container
}

func newSetChildFixture(t *testing.T) setChildFixture {
	t.Helper()

	var captured interface{}

	arraySchema := &schema.Array{Items: &schema.String{}}
	objectSchema := &schema.Object{
		ParameterNames: map[string]string{"Items": "items"},
		Properties: map[string]schema.Schema{
			"Items": arraySchema,
		},
	}

	var fixture setChildFixture
	fixture.setParamValue = &captured
	interpreter := &conflowfakes.FakeBlockInterpreter{}
	interpreter.SchemaReturns(objectSchema)
	blockInstance := &conflowfakes.FakeIdentifiable{}
	interpreter.CreateBlockReturns(blockInstance)
	interpreter.SetParamStub = func(_ conflow.Block, _ conflow.ID, value interface{}) error {
		*fixture.setParamValue = value
		return nil
	}

	blockNode := &conflowfakes.FakeBlockNode{}
	blockNode.IDReturns("parent")
	blockNode.InterpreterReturns(interpreter)

	logger := zerolog.NewDisabledLogger()
	evalCtx := conflow.NewEvalContext(context.Background(), nil, logger, nil, nil, blockNode)
	parentContainer := NewContainer(evalCtx, conflow.RuntimeConfig{}, blockNode, nil, nil, nil, false)
	parentContainer.block = blockInstance
	fixture.container = parentContainer

	return fixture
}

func setChildWithValue(parentContainer *Container, name conflow.ID, value interface{}) parsley.Error {
	arraySchema := &schema.Array{Items: &schema.String{}}
	paramNode := parameter.NewNode(
		"parent",
		conflow.NewIDNode(name, conflow.ClassifierNone, parsley.NilPos, parsley.NilPos),
		nil,
		true,
		nil,
	)
	paramNode.SetSchema(arraySchema)
	child := parameter.NewContainer(parentContainer.evalCtx, paramNode, parentContainer, value, nil, false)
	return parentContainer.setChild(child)
}
