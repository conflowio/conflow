// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package values

type Map[K comparable, V any] struct {
	m map[K]V
}

func NewMap[K comparable, V any](entries ...struct {
	Key   K
	Value V
}) *Map[K, V] {
	m := make(map[K]V, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return &Map[K, V]{m: m}
}

func NewMapFromGoMap[K comparable, V any](m map[K]V) *Map[K, V] {
	cp := make(map[K]V, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return &Map[K, V]{m: cp}
}

func MapOf[K comparable, V any](m map[K]V) *Map[K, V] {
	return NewMapFromGoMap(m)
}

func (m *Map[K, V]) Len() int { return len(m.m) }

func (m *Map[K, V]) Get(k K) (V, bool) {
	v, ok := m.m[k]
	return v, ok
}

func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}
	return keys
}

// GoMap returns a copy for interop (json, stdlib). Callers may mutate the copy.
func (m *Map[K, V]) GoMap() map[K]V {
	cp := make(map[K]V, len(m.m))
	for k, v := range m.m {
		cp[k] = v
	}
	return cp
}

type MapBuilder[K comparable, V any] struct {
	m map[K]V
}

func NewMapBuilder[K comparable, V any]() *MapBuilder[K, V] {
	return &MapBuilder[K, V]{m: make(map[K]V)}
}

func (b *MapBuilder[K, V]) Set(k K, v V) {
	b.m[k] = v
}

func (b *MapBuilder[K, V]) Freeze() *Map[K, V] {
	return NewMapFromGoMap(b.m)
}
