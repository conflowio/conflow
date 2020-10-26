// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import "time"

func BoolPtr(v bool) *bool {
	return &v
}

func BoolValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func Float32Ptr(v float32) *float32 {
	return &v
}

func Float32Value(v *float32) float32 {
	if v == nil {
		return 0
	}
	return *v
}

func Float64Ptr(v float64) *float64 {
	return &v
}

func Float64Value(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func IntPtr(v int) *int {
	return &v
}

func IntValue(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func Int8Ptr(v int8) *int8 {
	return &v
}

func Int8Value(v *int8) int8 {
	if v == nil {
		return 0
	}
	return *v
}

func Int16Ptr(v int16) *int16 {
	return &v
}

func Int16Value(v *int16) int16 {
	if v == nil {
		return 0
	}
	return *v
}

func Int32Ptr(v int32) *int32 {
	return &v
}

func Int32Value(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func Int64Ptr(v int64) *int64 {
	return &v
}

func Int64Value(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func StringPtr(v string) *string {
	return &v
}

func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func TimeDurationPtr(v time.Duration) *time.Duration {
	return &v
}

func TimeDurationValue(v *time.Duration) time.Duration {
	if v == nil {
		return 0
	}

	return *v
}

func UintPtr(v uint) *uint {
	return &v
}

func UintValue(v *uint) uint {
	if v == nil {
		return 0
	}
	return *v
}

func Uint8Ptr(v uint8) *uint8 {
	return &v
}

func Uint8Value(v *uint8) uint8 {
	if v == nil {
		return 0
	}
	return *v
}

func Uint16Ptr(v uint16) *uint16 {
	return &v
}

func Uint16Value(v *uint16) uint16 {
	if v == nil {
		return 0
	}
	return *v
}

func Uint32Ptr(v uint32) *uint32 {
	return &v
}

func Uint32Value(v *uint32) uint32 {
	if v == nil {
		return 0
	}
	return *v
}

func Uint64Ptr(v uint64) *uint64 {
	return &v
}

func Uint64Value(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}
