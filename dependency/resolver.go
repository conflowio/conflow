package dependency

import (
	"strings"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/basil/basil"
)

// Resolver will resolve all process dependencies in a workflow
// It will look for non-existing dependencies and cycles
//
// It uses Tarjan's strongly connected components algorithm to detect cycles.
type Resolver struct {
	nodes         map[basil.ID]*node
	providedNodes map[basil.ID]basil.ID
	index         int
	stack         stack
	result        []basil.Node
	dependencies  []basil.IdentifiableNode
}

// NewResolver creates a new dependency resolver
func NewResolver(nodes ...basil.Node) *Resolver {
	r := &Resolver{
		nodes:         make(map[basil.ID]*node),
		providedNodes: make(map[basil.ID]basil.ID),
	}
	r.AddNodes(nodes...)
	return r
}

// AddNode adds a new node to the dependency graph
func (r *Resolver) AddNodes(nodes ...basil.Node) {
	for _, n := range nodes {
		r.nodes[n.ID()] = &node{
			Node:  n,
			Index: -1,
		}
		for _, provided := range n.Provides() {
			r.providedNodes[provided] = n.ID()
		}
	}
}

// Resolve will resolve the dependency graph
func (r *Resolver) Resolve() (result []basil.Node, dependencies []basil.IdentifiableNode, err parsley.Error) {
	for _, v := range r.nodes {
		if v.Index == -1 {
			if err := r.strongConnect(v); err != nil {
				return nil, nil, err
			}
		}
	}
	return r.result, r.dependencies, nil
}

// strongConnect will find all the strongly connected components in the dependency graph based on Tarjan's algorithm
func (r *Resolver) strongConnect(v *node) parsley.Error {
	v.Index = r.index
	v.LowLink = r.index
	r.index++
	r.stack.Push(v)
	for _, d := range v.Node.Dependencies() {
		w, found := r.nodes[d.ID()]
		if !found {
			w, found = r.nodes[d.ParentID()]
		}
		if !found {
			if providerID, ok := r.providedNodes[d.ParentID()]; ok {
				w = r.nodes[providerID]
				found = true
			}
		}

		if !found {
			r.dependencies = append(r.dependencies, d)
			continue
		}

		if v.Node.ID() == w.Node.ID() {
			return parsley.NewErrorf(d.Pos(), "%s should not reference itself", d.ID())
		}

		if err := r.processEdge(v, w); err != nil {
			return err
		}
	}
	if v.LowLink == v.Index {
		if err := r.createComponent(v); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) processEdge(v, w *node) parsley.Error {
	if w.Index == -1 {
		if err := r.strongConnect(w); err != nil {
			return err
		}
		v.LowLink = min(v.LowLink, w.LowLink)
	} else if w.OnStack {
		v.LowLink = min(v.LowLink, w.Index)
	}

	return nil
}

func (r *Resolver) createComponent(v *node) parsley.Error {
	var component []basil.Node
	for {
		w := r.stack.Pop()
		component = append(component, w.Node)
		if w.Node.ID() == v.Node.ID() {
			break
		}
	}

	if len(component) > 1 {
		var ids []string
		for _, c := range component {
			ids = append(ids, string(c.ID()))
		}
		return parsley.NewErrorf(component[0].Pos(), "circular dependency detected: %s", strings.Join(ids, ", "))
	}

	r.result = append(r.result, component[0])

	return nil
}

func min(i1 int, i2 int) int {
	if i1 <= i2 {
		return i1
	}
	return i2
}
