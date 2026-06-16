// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package values

import "fmt"

type List[T any] struct {
	elems []T
}

func NewList[T any](elems ...T) *List[T] {
	cp := make([]T, len(elems))
	copy(cp, elems)
	return &List[T]{elems: cp}
}

func NewListFromSlice[T any](s []T) *List[T] {
	cp := make([]T, len(s))
	copy(cp, s)
	return &List[T]{elems: cp}
}

func ListOf[T any](elems ...T) *List[T] {
	return NewList(elems...)
}

func (l *List[T]) Len() int { return len(l.elems) }

func (l *List[T]) At(i int) T {
	if i < 0 || i >= len(l.elems) {
		panic(fmt.Sprintf("values.List.At: index %d out of range [0,%d)", i, len(l.elems)))
	}
	return l.elems[i]
}

// Elems returns a copy for interop (json, stdlib). Callers may mutate the copy.
func (l *List[T]) Elems() []T {
	cp := make([]T, len(l.elems))
	copy(cp, l.elems)
	return cp
}

type ListBuilder[T any] struct {
	elems []T
}

func NewListBuilder[T any]() *ListBuilder[T] {
	return &ListBuilder[T]{elems: make([]T, 0)}
}

func (b *ListBuilder[T]) Append(v T) {
	b.elems = append(b.elems, v)
}

func (b *ListBuilder[T]) Freeze() *List[T] {
	return NewListFromSlice(b.elems)
}
