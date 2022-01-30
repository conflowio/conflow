// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"strconv"

	"github.com/conflowio/conflow/src/conflow"
)

type idRegistry struct {
	ids    map[conflow.ID]struct{}
	nextID int
}

func newIDRegistry() *idRegistry {
	return &idRegistry{
		ids: map[conflow.ID]struct{}{},
	}
}

func (r *idRegistry) IDExists(id conflow.ID) bool {
	_, exists := r.ids[id]
	return exists
}

func (r *idRegistry) GenerateID() conflow.ID {
	id := conflow.ID(strconv.Itoa(r.nextID))
	r.ids[id] = struct{}{}
	r.nextID++
	return id
}

func (r *idRegistry) RegisterID(id conflow.ID) error {
	r.ids[id] = struct{}{}
	return nil
}
