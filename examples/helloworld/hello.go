// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/conflowio/conflow/pkg/conflow"
)

// Hello is capable to print some greetings
// @block "task"
type Hello struct {
	// @id
	id conflow.ID
	// @required
	to string
	// @read_only
	greeting string
	// @ignore
	r *rand.Rand
}

func (h *Hello) ID() conflow.ID {
	return h.id
}

// Init will initialise the random generator
func (h *Hello) Init(ctx context.Context) (bool, error) {
	h.r = rand.New(rand.NewSource(time.Now().Unix()))
	return false, nil
}

// Main will generate a random greeting
func (h *Hello) Run(ctx context.Context) (conflow.Result, error) {
	greetings := []string{"Hello", "Hi", "Hey", "Yo", "Sup"}

	h.greeting = fmt.Sprintf("%s %s!", greetings[h.r.Intn(len(greetings))], h.to)

	return nil, nil
}
