package generator_test

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/opsidian/basil/basil/function/generator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseArguments", func() {

	var source string
	var arguments []*generator.Argument
	var results []*generator.Argument
	var argumentsErr, resultsErr error

	JustBeforeEach(func() {
		file, err := parser.ParseFile(token.NewFileSet(), "testfile", source, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())

		str := file.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.FuncType)
		arguments, argumentsErr = generator.ParseArguments(str, file)
		results, resultsErr = generator.ParseResults(str, file)
	})

	Context("when the function has no arguments", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo() int64
			`
		})

		It("should return with empty argument list", func() {
			Expect(argumentsErr).ToNot(HaveOccurred())
			Expect(arguments).To(Equal([]*generator.Argument{}))
		})

		It("should return with the results", func() {
			Expect(resultsErr).ToNot(HaveOccurred())
			Expect(results).To(Equal([]*generator.Argument{
				{
					Name: "",
					Type: "int64",
				},
			}))
		})
	})

})
