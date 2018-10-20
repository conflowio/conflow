package ocl_test

import (
	"errors"

	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/ocl/function/functionfakes"
	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/ocl/variable"
	"github.com/opsidian/ocl/variable/variablefakes"
)

var _ = Describe("Context", func() {
	Describe("GetVar", func() {
		It("should call provider", func() {
			variableProvider := &variablefakes.FakeProvider{}
			variableProvider.GetVarReturns("RES", true)

			functionRegistry := &functionfakes.FakeRegistry{}

			ctx := ocl.NewContext(variableProvider, functionRegistry)
			res, found := ctx.GetVar("name")

			Expect(res).To(Equal("RES"))
			Expect(found).To(BeTrue())

			Expect(variableProvider.GetVarCallCount()).To(Equal(1))
			passedName := variableProvider.GetVarArgsForCall(0)
			Expect(passedName).To(Equal("name"))
		})
	})

	Describe("LookupVar", func() {
		It("should call provider", func() {
			err := errors.New("ERR")
			variableProvider := &variablefakes.FakeProvider{}
			variableProvider.LookupVarReturns("RES", err)

			functionRegistry := &functionfakes.FakeRegistry{}

			ctx := ocl.NewContext(variableProvider, functionRegistry)
			lookUp := func(provider variable.Provider) (interface{}, error) {
				return nil, nil // The return values don't matter here
			}

			res, lookupErr := ctx.LookupVar(lookUp)
			Expect(res).To(Equal("RES"))
			Expect(lookupErr).To(MatchError(err))

			Expect(variableProvider.LookupVarCallCount()).To(Equal(1))
		})
	})

	Describe("CallFunction", func() {
		It("should call registry", func() {
			res := "RES"
			err := parsley.NewError(parsley.Pos(1), errors.New("ERR"))
			variableProvider := &variablefakes.FakeProvider{}
			functionRegistry := &functionfakes.FakeRegistry{}
			functionRegistry.CallFunctionReturns(res, err)

			ctx := ocl.NewContext(variableProvider, functionRegistry)

			evalCtx := "CTX"
			function := &parsleyfakes.FakeNode{}
			function.ValueReturns("FUNCTION", nil)
			params := []parsley.Node{&parsleyfakes.FakeNode{}}

			actualRes, actualErr := ctx.CallFunction(evalCtx, function, params)
			Expect(actualRes).To(BeEquivalentTo(res))
			Expect(actualErr).To(BeEquivalentTo(err))

			Expect(functionRegistry.CallFunctionCallCount()).To(Equal(1))

			passedEvalCtx, passedFunction, passedParams := functionRegistry.CallFunctionArgsForCall(0)
			Expect(passedEvalCtx).To(BeEquivalentTo(evalCtx))
			Expect(passedFunction).To(BeEquivalentTo(function))
			Expect(passedParams).To(BeEquivalentTo(params))
		})
	})

	Describe("FunctionExists", func() {
		It("should call registry", func() {
			variableProvider := &variablefakes.FakeProvider{}
			functionRegistry := &functionfakes.FakeRegistry{}
			functionRegistry.FunctionExistsReturnsOnCall(0, false)
			functionRegistry.FunctionExistsReturnsOnCall(1, true)

			ctx := ocl.NewContext(variableProvider, functionRegistry)

			res := ctx.FunctionExists("FOO")
			Expect(res).To(BeFalse())

			res = ctx.FunctionExists("BAR")
			Expect(res).To(BeTrue())

			Expect(functionRegistry.FunctionExistsCallCount()).To(Equal(2))

			passedName := functionRegistry.FunctionExistsArgsForCall(0)
			Expect(passedName).To(Equal("FOO"))
		})
	})

	Describe("RegisterFunction", func() {
		It("should call registry", func() {
			variableProvider := &variablefakes.FakeProvider{}
			functionRegistry := &functionfakes.FakeRegistry{}

			ctx := ocl.NewContext(variableProvider, functionRegistry)

			f := &functionfakes.FakeCallable{}
			f.CallFunctionReturns("FOO", nil)

			ctx.RegisterFunction("FOO", f)

			Expect(functionRegistry.RegisterFunctionCallCount()).To(Equal(1))

			passedName, passedFunction := functionRegistry.RegisterFunctionArgsForCall(0)
			Expect(passedName).To(BeEquivalentTo("FOO"))
			Expect(passedFunction).To(BeEquivalentTo(f))
		})
	})
})
