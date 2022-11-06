// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

// Epsilon is used as a float64 comparison tolerance
const Epsilon = 0.000000001

type Comparator[T any] func(v1, v2 T) bool

func SelfComparator[T interface{ Equal(v T) bool }]() Comparator[T] {
	return func(v1, v2 T) bool {
		return v1.Equal(v2)
	}
}

func ArrayEquals[T any](c Comparator[T]) Comparator[[]T] {
	return func(v1 []T, v2 []T) bool {
		if len(v1) != len(v2) {
			return false
		}
		for i, e1 := range v1 {
			if !c(e1, v2[i]) {
				return false
			}
		}
		return true
	}
}

func MapEquals[T any](c Comparator[T]) Comparator[map[string]T] {
	return func(v1 map[string]T, v2 map[string]T) bool {
		if len(v1) != len(v2) {
			return false
		}
		for k, e1 := range v1 {
			e2, ok := v2[k]
			if !ok {
				return false
			}
			if !c(e1, e2) {
				return false
			}
		}
		return true
	}
}

func PointerEquals[T any](c Comparator[T]) Comparator[*T] {
	return func(v1, v2 *T) bool {
		switch {
		case v1 == nil && v2 == nil:
			return true
		case v1 == nil || v2 == nil:
			return false
		default:
			return c(*v1, *v2)
		}
	}
}

func ValueEquals[T ~bool | ~int64 | ~string](v1, v2 T) bool {
	return v1 == v2
}

func FloatEquals[T ~float64](v1, v2 T) bool {
	return v1-v2 < Epsilon && v2-v1 < Epsilon
}
