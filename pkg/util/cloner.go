// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import "time"

type Cloner[T any] func(v T) T

func SelfCloner[T interface{ Clone() T }]() Cloner[T] {
	return func(v T) T {
		return v.Clone()
	}
}

func CloneArray[T any](c Cloner[T]) Cloner[[]T] {
	return func(v []T) []T {
		if v == nil {
			return nil
		}
		res := make([]T, 0, len(v))
		for _, e := range v {
			res = append(res, c(e))
		}
		return res
	}
}

func CloneMap[T any](c Cloner[T]) Cloner[map[string]T] {
	return func(v map[string]T) map[string]T {
		if v == nil {
			return nil
		}
		res := make(map[string]T, len(v))
		for k, e := range v {
			res[k] = c(e)
		}
		return res
	}
}

func ClonePointer[T any](c Cloner[T]) Cloner[*T] {
	return func(v *T) *T {
		if v == nil {
			return nil
		}

		res := c(*v)
		return &res
	}
}

func CloneValue[T ~bool | ~int64 | ~float64 | ~string | time.Time](v T) T {
	return v
}
