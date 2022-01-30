// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

import (
	"errors"
	"fmt"
	"sync"
)

var registry = Registry{
	schemas: map[string]interface{}{},
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
	schemas map[string]interface{}
}

func (r *Registry) RegisterSchema(s Schema) {
	if s.GetID() == "" {
		panic(errors.New("schema can not be registered without an id"))
	}

	registryLock.Lock()
	r.schemas[s.GetID()] = s
	registryLock.Unlock()
}

func (r *Registry) RegisterResolver(uri string, res Resolver) {
	registryLock.Lock()
	r.schemas[uri] = res
	registryLock.Unlock()
}

func (r *Registry) GetSchema(uri string) (Schema, error) {
	registryLock.RLock()
	s, found := r.schemas[uri]
	registryLock.RUnlock()

	if !found {
		return nil, nil
	}

	switch st := s.(type) {
	case nil:
		return nil, nil
	case Schema:
		return st, nil
	case Resolver:
		registryLock.Lock()
		defer registryLock.Unlock()

		// Re-check the value of s, in case it was already resolved by another GetSchema acquiring the lock first
		if _, ok := r.schemas[uri].(Resolver); !ok {
			return Get(uri)
		}

		s, err := st.ResolveSchema(uri)
		if err != nil {
			return nil, err
		}

		r.schemas[uri] = s

		return s, nil
	default:
		panic(fmt.Errorf("unexpected schema entry: %T", s))
	}
}
