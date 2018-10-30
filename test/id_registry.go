package test

import "strconv"

type idRegistry struct {
	ids    map[string]struct{}
	nextID int
}

func newIDRegistry() *idRegistry {
	return &idRegistry{
		ids: map[string]struct{}{},
	}
}

func (r *idRegistry) IDExists(id string) bool {
	_, exists := r.ids[id]
	return exists
}

func (r *idRegistry) GenerateID() string {
	id := strconv.Itoa(r.nextID)
	r.ids[id] = struct{}{}
	r.nextID++
	return id
}

func (r *idRegistry) RegisterID(id string) error {
	r.ids[id] = struct{}{}
	return nil
}
