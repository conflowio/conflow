// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package variable_test

import (
	"context"

	"github.com/conflowio/parsley/parsley"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/conflowfakes"
	"github.com/conflowio/conflow/pkg/conflow/variable"
	"github.com/conflowio/conflow/pkg/loggers/zerolog"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/values"
)

var _ = Describe("Isolation", func() {
	var evalCtx *conflow.EvalContext
	var blockContainer *conflowfakes.FakeBlockContainer
	var blockNode *conflowfakes.FakeBlockNode
	var node *variable.Node

	arraySchema := &schema.Array{Items: &schema.String{}}

	BeforeEach(func() {
		blockNode = &conflowfakes.FakeBlockNode{}
		blockNode.IDReturns("foo")
		blockNode.GetPropertySchemaStub = func(id conflow.ID) (schema.Schema, bool) {
			if string(id) == "items" {
				return arraySchema, true
			}
			return nil, false
		}

		parseCtx := conflow.NewParseContext(nil, conflow.NewIDRegistry(8, 16), nil)
		Expect(parseCtx.AddBlockNode(blockNode)).To(Succeed())

		blockContainer = &conflowfakes.FakeBlockContainer{}
		blockContainer.NodeReturns(blockNode)

		logger := zerolog.NewDisabledLogger()
		evalCtx = conflow.NewEvalContext(
			context.Background(),
			nil,
			logger,
			nil,
			map[conflow.ID]conflow.BlockContainer{"foo": blockContainer},
			blockNode,
		)

		node = variable.NewNode(
			conflow.NewIDNode("foo", conflow.ClassifierNone, parsley.NilPos, parsley.NilPos),
			conflow.NewIDNode("items", conflow.ClassifierNone, parsley.NilPos, parsley.NilPos),
		)
		Expect(node.StaticCheck(parseCtx)).To(Succeed())
	})

	It("returns a bound copy for array parameters", func() {
		upstream := []interface{}{"shared"}
		blockContainer.ParamReturns(upstream)

		value, err := node.Value(evalCtx)
		Expect(err).NotTo(HaveOccurred())

		list, ok := value.(*values.List[interface{}])
		Expect(ok).To(BeTrue())
		Expect(list.At(0)).To(Equal("shared"))
		Expect(value).NotTo(BeIdenticalTo(&upstream))

		upstream[0] = "mutated"
		Expect(list.At(0)).To(Equal("shared"))
	})

	Context("with a string parameter schema", func() {
		BeforeEach(func() {
			blockNode.GetPropertySchemaStub = func(id conflow.ID) (schema.Schema, bool) {
				if string(id) == "param1" {
					return schema.StringValue(), true
				}
				return nil, false
			}

			parseCtx := conflow.NewParseContext(nil, conflow.NewIDRegistry(8, 16), nil)
			Expect(parseCtx.AddBlockNode(blockNode)).To(Succeed())

			node = variable.NewNode(
				conflow.NewIDNode("foo", conflow.ClassifierNone, parsley.NilPos, parsley.NilPos),
				conflow.NewIDNode("param1", conflow.ClassifierNone, parsley.NilPos, parsley.NilPos),
			)
			Expect(node.StaticCheck(parseCtx)).To(Succeed())
		})

		It("passes string parameters through unchanged", func() {
			blockContainer.ParamReturns("bar")

			value, err := node.Value(evalCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("bar"))
		})
	})
})
