// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow_test

import (
	"context"

	"github.com/conflowio/parsley/parsley"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/conflowfakes"
	"github.com/conflowio/conflow/pkg/conflow/job"
	"github.com/conflowio/conflow/pkg/loggers/zerolog"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/values"

	_ "github.com/conflowio/conflow/pkg/test"
)

var _ = Describe("Evaluate input isolation", func() {
	var parseCtx *conflow.ParseContext
	var blockNode *conflowfakes.FakeBlockNode
	var boundInput interface{}
	var upstream []interface{}
	var inputParams map[conflow.ID]interface{}

	BeforeEach(func() {
		boundInput = nil
		upstream = []interface{}{"shared"}

		baseSchema, err := schema.Get("github.com/conflowio/conflow/pkg/test.Block")
		Expect(err).NotTo(HaveOccurred())
		objectSchema := baseSchema.Copy().(*schema.Object)
		objectSchema.Properties["FieldArray"].(schema.MetadataAccessor).SetAnnotation(annotations.UserDefined, "true")

		interpreter := &conflowfakes.FakeBlockInterpreter{}
		interpreter.SchemaReturns(objectSchema)

		blockNode = &conflowfakes.FakeBlockNode{}
		blockNode.IDReturns("foo")
		blockNode.InterpreterReturns(interpreter)
		blockNode.ValueStub = func(ctx interface{}) (interface{}, parsley.Error) {
			boundInput = ctx.(*conflow.EvalContext).InputParams["field_array"]
			return &conflowfakes.FakeIdentifiable{}, nil
		}

		parseCtx = conflow.NewParseContext(nil, conflow.NewIDRegistry(8, 16), nil)
		Expect(parseCtx.AddBlockNode(blockNode)).To(Succeed())

		inputParams = map[conflow.ID]interface{}{
			"field_array": upstream,
		}
	})

	It("binds array @input parameters before validation", func() {
		logger := zerolog.NewDisabledLogger()
		scheduler := job.NewScheduler(logger, 1, 10)
		scheduler.Start()
		defer scheduler.Stop()

		_, err := conflow.Evaluate(
			parseCtx,
			context.Background(),
			nil,
			logger,
			scheduler,
			"foo",
			inputParams,
		)
		Expect(err).NotTo(HaveOccurred())

		list, ok := boundInput.(*values.List[interface{}])
		Expect(ok).To(BeTrue())
		Expect(list.At(0)).To(Equal("shared"))

		upstream[0] = "mutated"
		Expect(list.At(0)).To(Equal("shared"))
	})

	It("preserves immutable list @input without slice round-trip", func() {
		upstream := values.ListOf[interface{}]("shared")
		inputParams = map[conflow.ID]interface{}{
			"field_array": upstream,
		}

		logger := zerolog.NewDisabledLogger()
		scheduler := job.NewScheduler(logger, 1, 10)
		scheduler.Start()
		defer scheduler.Stop()

		_, err := conflow.Evaluate(
			parseCtx,
			context.Background(),
			nil,
			logger,
			scheduler,
			"foo",
			inputParams,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(boundInput).To(BeIdenticalTo(upstream))
	})
})
