package ocl_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/ocl/variable"
	"github.com/opsidian/ocl/variable/variablefakes"
)

var _ = Describe("Context", func() {
	Describe("GetVar", func() {
		It("should call provider", func() {
			variableProvider := &variablefakes.FakeProvider{}
			variableProvider.GetVarReturns("RES", true)

			ctx := ocl.NewContext(variableProvider)
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

			ctx := ocl.NewContext(variableProvider)
			lookUp := func(provider variable.Provider) (interface{}, error) {
				return nil, nil // The return values don't matter here
			}

			res, lookupErr := ctx.LookupVar(lookUp)
			Expect(res).To(Equal("RES"))
			Expect(lookupErr).To(MatchError(err))

			Expect(variableProvider.LookupVarCallCount()).To(Equal(1))
		})
	})
})

// func TestContextLookupVarShouldCallProvider(t *testing.T) {
// 	err := errors.New("ERR")
// 	variableProvider := new(mocks.Provider)
// 	variableProvider.On("LookupVar", mock.AnythingOfType("variable.LookUp")).Return("RES", err)

// 	ctx := flint.NewContext(variableProvider)
// 	lookUp := func(provider variable.Provider) (interface{}, error) {
// 		return nil, nil // The returns values don't matter here
// 	}
// 	res, expErr := ctx.LookupVar(lookUp)

// 	variableProvider.AssertExpectations(t)
// 	assert.Equal(t, res, "RES")
// 	assert.Equal(t, expErr, err)
// }
