package dependency_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/basilfakes"
	"github.com/opsidian/basil/basil/dependency"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Resolver", func() {

	var resolver *dependency.Resolver
	var result []basil.Node
	var dependencies basil.Dependencies
	var err error

	BeforeEach(func() {
		resolver = dependency.NewResolver("b")
	})

	JustBeforeEach(func() {
		result, dependencies, err = resolver.Resolve()
	})

	Context("when there are no nodes", func() {
		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with empty result", func() {
			Expect(result).To(BeNil())
		})

		It("should return with no dependencies", func() {
			Expect(dependencies).To(BeNil())
		})
	})

	Context("when there is only one node", func() {
		var param1 *basilfakes.FakeNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			resolver.AddNodes(param1)
		})

		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with the node", func() {
			Expect(result).To(ConsistOf(param1))
		})

		It("should return with no parent dependencies", func() {
			Expect(dependencies).To(BeNil())
		})
	})

	Context("when the nodes don't have dependencies", func() {
		var param1, param2 *basilfakes.FakeNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			param2 = &basilfakes.FakeNode{}
			param2.IDReturns("b.param2")
			resolver.AddNodes(param1, param2)
		})

		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with all nodes in any order", func() {
			Expect(result).To(ConsistOf(param1, param2))
		})

		It("should return with no dependencies", func() {
			Expect(dependencies).To(BeNil())
		})
	})

	Context("when the nodes have dependencies", func() {
		var param1, param2, param3, param4 *basilfakes.FakeNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			param2 = &basilfakes.FakeNode{}
			param2.IDReturns("b.param2")
			param3 = &basilfakes.FakeNode{}
			param3.IDReturns("b.param3")
			param4 = &basilfakes.FakeNode{}
			param4.IDReturns("b.param4")

			param1.DependenciesReturns(basil.Dependencies{"b.param2": dep("b.param2")})
			param3.DependenciesReturns(basil.Dependencies{"b.param4": dep("b.param4")})

			resolver.AddNodes(param1, param2, param3, param4)
		})

		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with no dependencies", func() {
			Expect(dependencies).To(BeNil())
		})

		It("param2 should be resolved before param1", func() {
			Expect(nodeIndex(result, param2) < nodeIndex(result, param1)).To(BeTrue())
		})

		It("param4 should be resolved before param3", func() {
			Expect(nodeIndex(result, param4) < nodeIndex(result, param3)).To(BeTrue())
		})
	})

	Context("when a child node is the dependency", func() {
		var param1, node2 *basilfakes.FakeNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			node2 = &basilfakes.FakeNode{}
			node2.IDReturns("node2")

			param1.DependenciesReturns(basil.Dependencies{"node2.x": dep("node2.x")})

			resolver.AddNodes(param1, node2)
		})

		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with no dependencies", func() {
			Expect(dependencies).To(BeNil())
		})

		It("node2 should be resolved before param1", func() {
			Expect(nodeIndex(result, node2) < nodeIndex(result, param1)).To(BeTrue())
		})
	})

	Context("when the nodes have transitive dependencies", func() {
		var param1, node2 *basilfakes.FakeNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			node2 = &basilfakes.FakeNode{}
			node2.IDReturns("node2")

			node2.ProvidesReturns([]basil.ID{"node3"})
			param1.DependenciesReturns(basil.Dependencies{"node3.x": dep("node3.x")})

			resolver.AddNodes(param1, node2)
		})

		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with no dependencies", func() {
			Expect(dependencies).To(BeNil())
		})

		It("node2 should be resolved before param1", func() {
			Expect(nodeIndex(result, node2) < nodeIndex(result, param1)).To(BeTrue())
		})
	})

	Context("when the nodes have circular dependencies", func() {
		var param1, param2, param3 *basilfakes.FakeNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			param1.PosReturns(parsley.Pos(1))
			param2 = &basilfakes.FakeNode{}
			param2.IDReturns("b.param2")
			param2.PosReturns(parsley.Pos(2))
			param3 = &basilfakes.FakeNode{}
			param3.IDReturns("b.param3")
			param3.PosReturns(parsley.Pos(3))

			param1.DependenciesReturns(basil.Dependencies{"b.param2": dep("b.param2")})
			param2.DependenciesReturns(basil.Dependencies{"b.param1": dep("b.param1")})

			resolver.AddNodes(param1, param2, param3)
		})

		It("should return with an error", func() {
			err1 := parsley.NewError(parsley.Pos(2), errors.New("circular dependency detected: b.param2, b.param1"))
			err2 := parsley.NewError(parsley.Pos(1), errors.New("circular dependency detected: b.param1, b.param2"))
			Expect(err).To(Or(MatchError(err1), MatchError(err2)))
		})
	})

	Context("when a node is referencing itself", func() {
		var param1, param2 *basilfakes.FakeNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			param1.PosReturns(parsley.Pos(1))
			param2 = &basilfakes.FakeNode{}
			param2.IDReturns("b.param2")
			param2.PosReturns(parsley.Pos(2))

			param1.DependenciesReturns(basil.Dependencies{"b.param2": dep("b.param2")})
			param2.DependenciesReturns(basil.Dependencies{"b.param2": dep("b.param2")})

			resolver.AddNodes(param1, param2)
		})

		It("should return with an error", func() {
			err := parsley.NewError(parsley.Pos(1), errors.New("b.param1 should not reference itself"))
			Expect(err).To(MatchError(err))
		})
	})

	Context("when a node is referencing an unknown parameter", func() {
		var param1 *basilfakes.FakeNode
		var dep1 *basilfakes.FakeVariableNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			param1.PosReturns(parsley.Pos(1))

			dep1 = dep("b.param2")
			dep1.PosReturns(parsley.Pos(2))

			param1.DependenciesReturns(basil.Dependencies{"b.param2": dep1})

			resolver.AddNodes(param1)
		})

		It("should return with an error", func() {
			expectedErr := parsley.NewError(parsley.Pos(2), errors.New("unknown parameter: \"b.param2\""))
			Expect(err).To(MatchError(expectedErr))
		})

	})

	Context("when a node is referencing an external parameter", func() {
		var param1 *basilfakes.FakeNode
		var dep1 *basilfakes.FakeVariableNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			param1.PosReturns(parsley.Pos(1))

			dep1 = dep("b2.param2")

			param1.DependenciesReturns(basil.Dependencies{"b2.param2": dep1})

			resolver.AddNodes(param1)
		})

		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with all nodes in any order", func() {
			Expect(result).To(ConsistOf(param1))
		})

		It("should return the node as a dependency", func() {
			Expect(dependencies).To(ConsistOf(dep1))
		})
	})

	Context("when a node is referencing an unknown node", func() {
		var param1 *basilfakes.FakeNode
		var dep1 *basilfakes.FakeVariableNode

		BeforeEach(func() {
			param1 = &basilfakes.FakeNode{}
			param1.IDReturns("b.param1")
			param1.PosReturns(parsley.Pos(1))

			dep1 = dep("b2.param2")

			param1.DependenciesReturns(basil.Dependencies{"b2.param2": dep1})

			resolver.AddNodes(param1)
		})

		It("should return with no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with all nodes in any order", func() {
			Expect(result).To(ConsistOf(param1))
		})

		It("should return the node as a dependency", func() {
			Expect(dependencies).To(ConsistOf(dep1))
		})
	})
})

func dep(id string) *basilfakes.FakeVariableNode {
	f := &basilfakes.FakeVariableNode{}
	f.IDReturns(basil.ID(id))
	f.ParentIDReturns(basil.ID(id[0:strings.IndexByte(id, '.')]))
	return f
}

func nodeIndex(l []basil.Node, n basil.Node) int {
	for i, ln := range l {
		if ln == n {
			return i
		}
	}
	return -1
}
