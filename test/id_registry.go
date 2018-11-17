package test

import (
	"strconv"

	"github.com/opsidian/basil/basil"
)

type idRegistry struct {
	ids    map[basil.ID]struct{}
	nextID int
}

func newIDRegistry() *idRegistry {
	return &idRegistry{
		ids: map[basil.ID]struct{}{},
	}
}

func (r *idRegistry) IDExists(id basil.ID) bool {
	_, exists := r.ids[id]
	return exists
}

func (r *idRegistry) GenerateID() basil.ID {
	id := basil.ID(strconv.Itoa(r.nextID))
	r.ids[id] = struct{}{}
	r.nextID++
	return id
}

func (r *idRegistry) RegisterID(id basil.ID) error {
	r.ids[id] = struct{}{}
	return nil
}
