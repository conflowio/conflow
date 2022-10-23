// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package dependency

import (
	"errors"
	"strings"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
)

// Resolver will resolve all process dependencies in a workflow
// It will look for non-existing dependencies and cycles
//
// It uses Tarjan's strongly connected components algorithm to detect cycles.
type Resolver struct {
	id             conflow.ID
	nodes          map[conflow.ID]*node
	providedNodes  map[conflow.ID]conflow.ID
	generatedNodes map[conflow.ID]conflow.ID
	index          int
	stack          stack
	result         []conflow.Node
	dependencies   conflow.Dependencies
}

// NewResolver creates a new dependency resolver
func NewResolver(id conflow.ID, nodes ...conflow.Node) *Resolver {
	r := &Resolver{
		id:             id,
		nodes:          make(map[conflow.ID]*node),
		providedNodes:  make(map[conflow.ID]conflow.ID),
		generatedNodes: make(map[conflow.ID]conflow.ID),
	}
	r.AddNodes(nodes...)
	return r
}

// AddNode adds a new node to the dependency graph
func (r *Resolver) AddNodes(nodes ...conflow.Node) {
	for _, n := range nodes {
		r.nodes[n.ID()] = &node{
			Node:  n,
			Index: -1,
		}
		for _, id := range n.Provides() {
			r.providedNodes[id] = n.ID()
		}

		// If the node is a generator then we create an extra node so we have separate nodes for start and finish
		// Nodes depending on any of the generated nodes should depend on the start node
		// Nodes referencing any fields on the generator node should depend on the finish (original) node
		// We need to do this to avoid circular dependencies
		if len(n.Generates()) > 0 {
			for _, id := range n.Generates() {
				r.generatedNodes[id] = n.ID()
			}
			r.nodes[n.ID()+"-start"] = &node{
				Node:  n,
				Index: -1,
			}
			r.nodes[n.ID()].skip = true
			r.nodes[n.ID()].extraDependencies = []conflow.ID{n.ID() + "-start"}
		}
	}
}

// Resolve will resolve the dependency graph
func (r *Resolver) Resolve() (result []conflow.Node, dependencies conflow.Dependencies, err parsley.Error) {
	// We want to detect if a node depends on a generator node and any of its generated blocks
	// In this case we should return with a circular dependency error
	if len(r.generatedNodes) > 0 {
		for _, v := range r.nodes {
			if v.Index == -1 {
				r.generatorDependencies(v)
			}
		}
		for id, v := range r.nodes {
			v.Index = -1 // reset the Index for strongConnect
			for _, generatorID := range v.generatorDependencies {
				if _, isParam := v.Node.(conflow.ParameterNode); isParam {
					return nil, nil, parsley.NewError(v.Node.Pos(), errors.New("a parameter can not depend on a node generated in the same block"))
				}
				r.nodes[generatorID].extraDependencies = append(r.nodes[generatorID].extraDependencies, id)
			}
		}
	}

	for _, v := range r.nodes {
		if v.Index == -1 {
			if err := r.strongConnect(v); err != nil {
				return nil, nil, err
			}
		}
	}
	return r.result, r.dependencies, nil
}

func (r *Resolver) generatorDependencies(v *node) []conflow.ID {
	if v.Index == 0 {
		return v.generatorDependencies
	}

	var res []conflow.ID
	for _, d := range v.Node.Dependencies() {
		if generatorID, ok := r.generatedNodes[d.ParentID()]; ok {
			res = append(res, generatorID)
		}
		w, found := r.nodes[d.ParentID()]
		if !found {
			if providerID, ok := r.providedNodes[d.ParentID()]; ok {
				w = r.nodes[providerID]
				found = true
			}
		}
		if found {
			res = append(res, r.generatorDependencies(w)...)
		}
	}
	v.Index = 0
	v.generatorDependencies = res
	return res
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
			if d.ParentID() == r.id {
				return parsley.NewErrorf(d.Pos(), "unknown parameter: %q", d.ID())
			}

			w, found = r.nodes[d.ParentID()]
		}
		if !found {
			if providerID, ok := r.providedNodes[d.ParentID()]; ok {
				w = r.nodes[providerID]
				found = true
			}
		}

		if !found {
			if generatorID, ok := r.generatedNodes[d.ParentID()]; ok {
				w = r.nodes[generatorID+"-start"]
				found = true
			}
		}

		if !found {
			if r.dependencies == nil {
				r.dependencies = make(conflow.Dependencies)
			}
			r.dependencies[d.ID()] = d
			continue
		}

		if v.Node.ID() == w.Node.ID() {
			return parsley.NewErrorf(d.Pos(), "%s should not reference itself", d.ID())
		}

		if err := r.processEdge(v, w); err != nil {
			return err
		}
	}

	for _, id := range v.extraDependencies {
		w := r.nodes[id]
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
	var component []conflow.Node
	var hasSkipped bool
	for {
		w := r.stack.Pop()
		component = append(component, w.Node)
		if w.skip {
			hasSkipped = true
		}
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

	if !hasSkipped {
		r.result = append(r.result, component[0])
	}

	return nil
}

func min(i1 int, i2 int) int {
	if i1 <= i2 {
		return i1
	}
	return i2
}
