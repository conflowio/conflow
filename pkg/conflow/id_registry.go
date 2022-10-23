// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"fmt"
	"sync"

	"github.com/conflowio/conflow/pkg/util"
)

// IDRegistry generates and stores identifiers
type IDRegistry interface {
	IDExists(id ID) bool
	RegisterID(id ID) error
	GenerateID() ID
}

// idRegistry stores and generates identifiers
type idRegistry struct {
	ids       map[ID]struct{}
	minLength int
	maxLength int
	lock      sync.RWMutex
}

func NewIDRegistry(minLength int, maxLength int) IDRegistry {
	return &idRegistry{
		ids:       map[ID]struct{}{},
		minLength: minLength,
		maxLength: maxLength,
	}
}

// IDExists returns true if the identifier already exists
func (i *idRegistry) IDExists(id ID) bool {
	i.lock.RLock()
	defer i.lock.RUnlock()

	_, exists := i.ids[id]
	return exists
}

// GenerateID generates and registers a new block id
// If there are many block ids generated with the current minimum length and it's getting hard to generate unique ones
// then the min length will be increased by one (up to the maximum length)
func (i *idRegistry) GenerateID() ID {
	util.SeedMathRand()

	tries := 0
	for {
		id := ID("0x" + util.RandHexString(i.minLength, true))
		err := i.RegisterID(id)
		if err == nil {
			return id
		}
		tries++
		if tries == 3 {
			if i.minLength < i.maxLength {
				i.minLength++
				tries = 0
			} else {
				panic("unable to generate unique id, please increase the maximum identifier length")
			}
		}
	}
}

// RegisterID registers a new block id and returns an error if it is already taken
func (i *idRegistry) RegisterID(id ID) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	_, exists := i.ids[id]
	if exists {
		return fmt.Errorf("%q identifier already exists", id)
	}

	i.ids[id] = struct{}{}

	return nil
}
