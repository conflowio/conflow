// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"errors"
	"strings"
	"sync"
)

var registry = Registry{
	schemas:   map[string]Schema{},
	resolvers: map[string]Resolver{},
}
var registryLock = sync.RWMutex{}

func Get(uri string) (Schema, error) {
	return registry.GetSchema(uri)
}

func Register(s Schema) {
	registry.RegisterSchema(s)
}

func RegisterResolver(uri string, res Resolver) {
	registry.RegisterResolver(uri, res)
}

type Registry struct {
	schemas   map[string]Schema
	resolvers map[string]Resolver
}

func (r *Registry) RegisterSchema(s Schema) {
	if s.GetID() == "" {
		panic(errors.New("schema can not be registered without an id"))
	}

	registryLock.Lock()
	r.schemas[s.GetID()] = s
	registryLock.Unlock()
}

func (r *Registry) RegisterResolver(prefix string, res Resolver) {
	registryLock.Lock()
	r.resolvers[prefix] = res
	registryLock.Unlock()
}

func (r *Registry) GetSchema(uri string) (Schema, error) {
	registryLock.RLock()
	s, found := r.schemas[uri]
	registryLock.RUnlock()

	if found {
		return s, nil
	}

	registryLock.Lock()
	defer registryLock.Unlock()

	// Re-check the registry in case it was already resolved by another GetSchema acquiring the lock first
	if s, found := r.schemas[uri]; found {
		return s, nil
	}

	for prefix, resolver := range r.resolvers {
		if strings.HasPrefix(uri, prefix) {
			s, err := resolver.ResolveSchema(uri)
			if err != nil {
				return nil, err
			}

			if s == nil {
				continue
			}

			r.schemas[uri] = s

			return s, nil
		}
	}

	r.schemas[uri] = nil

	return nil, nil
}
