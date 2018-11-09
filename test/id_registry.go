package test

import (
	"strconv"

	"github.com/opsidian/basil/variable"
)

type idRegistry struct {
	ids    map[variable.ID]struct{}
	nextID int
}

func newIDRegistry() *idRegistry {
	return &idRegistry{
		ids: map[variable.ID]struct{}{},
	}
}

func (r *idRegistry) IDExists(id variable.ID) bool {
	_, exists := r.ids[id]
	return exists
}

func (r *idRegistry) GenerateID() variable.ID {
	id := variable.ID(strconv.Itoa(r.nextID))
	r.ids[id] = struct{}{}
	r.nextID++
	return id
}

func (r *idRegistry) RegisterID(id variable.ID) error {
	r.ids[id] = struct{}{}
	return nil
}
