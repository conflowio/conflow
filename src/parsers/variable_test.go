// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers_test

import (
	"context"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/schema"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/conflow/conflowfakes"
	"github.com/conflowio/conflow/src/loggers/zerolog"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/test"
)

var _ = Describe("Variable", func() {

	var p = parsers.Variable()
	var parsleyContext *parsley.Context
	var evalCtx *conflow.EvalContext
	var res parsley.Node
	var parseErr, evalErr error
	var value interface{}
	var input string
	var blockNode *conflowfakes.FakeBlockNode

	BeforeEach(func() {
		logger := zerolog.NewDisabledLogger()
		evalCtx = conflow.NewEvalContext(context.Background(), nil, logger, &test.Scheduler{}, nil)
		parseErr = nil
		evalErr = nil
		value = nil
		blockNode = nil
	})

	JustBeforeEach(func() {
		parsleyContext = test.ParseCtx(input, nil, nil)
		parseCtx := conflow.NewParseContext(
			parsleyContext.FileSet(),
			conflow.NewIDRegistry(8, 16),
			nil,
		)
		parsleyContext.SetUserContext(parseCtx)

		if blockNode != nil {
			err := parseCtx.AddBlockNode(blockNode)
			Expect(err).ToNot(HaveOccurred())
		}

		res, parseErr = parsley.Parse(parsleyContext, combinator.Sentence(p))
		if parseErr == nil {
			value, evalErr = parsley.EvaluateNode(evalCtx, res)
		}
	})

	Context("when referencing only the block name", func() {
		BeforeEach(func() {
			input = "foo"
		})

		It("should return with a parse error", func() {
			Expect(parseErr).To(HaveOccurred())
		})
	})

	Context("when referencing a block module parameter", func() {
		BeforeEach(func() {
			blockNode = &conflowfakes.FakeBlockNode{}
			blockNode.IDReturns(conflow.ID("foo"))
			blockNode.GetPropertySchemaStub = func(id conflow.ID) (schema.Schema, bool) {
				if string(id) == "param1" {
					return schema.StringValue(), true
				}
				return nil, false
			}

			cont := &conflowfakes.FakeBlockContainer{}
			cont.ParamReturnsOnCall(0, "bar")
			cont.NodeReturns(blockNode)

			ctx, cancel := context.WithCancel(context.Background())
			evalCtx = evalCtx.New(ctx, cancel, map[conflow.ID]conflow.BlockContainer{"foo": cont})
		})

		Context("with an existing parameter", func() {
			BeforeEach(func() {
				input = "foo.param1"
			})

			It("should evaluate successfully", func() {
				Expect(parseErr).ToNot(HaveOccurred())
				Expect(evalErr).ToNot(HaveOccurred())
				Expect(value).To(Equal("bar"))
			})
		})

		Context("with a nonexisting parameter", func() {
			BeforeEach(func() {
				input = "foo.param2"
			})

			It("should return a parse error", func() {
				Expect(parseErr).To(MatchError("parameter \"param2\" does not exist at testfile:1:5"))
			})
		})
	})

	Context("when referencing a non-existing block", func() {
		BeforeEach(func() {
			input = "foo.param1"
		})

		It("should return a parse error", func() {
			Expect(parseErr).To(MatchError("block \"foo\" does not exist at testfile:1:1"))
		})
	})

})
