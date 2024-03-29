// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package zerolog

import (
	"net"
	"time"

	"github.com/rs/zerolog"

	"github.com/conflowio/conflow/pkg/conflow"
)

type Array struct {
	arr *zerolog.Array
}

func (a *Array) MarshalZerologArray(*zerolog.Array) {
}

func (a *Array) MarshalLogArray(conflow.LogArray) {
}

func (a *Array) Object(obj conflow.LogObjectMarshaler) conflow.LogArray {
	a.arr = a.arr.Object(&ObjectMarshalerWrapper{obj: obj})
	return a
}

func (a *Array) ID(val conflow.ID) conflow.LogArray {
	a.arr = a.arr.Str(string(val))
	return a
}

func (a *Array) Str(val string) conflow.LogArray {
	a.arr = a.arr.Str(val)
	return a
}

func (a *Array) Bytes(val []byte) conflow.LogArray {
	a.arr = a.arr.Bytes(val)
	return a
}

func (a *Array) Hex(val []byte) conflow.LogArray {
	a.arr = a.arr.Hex(val)
	return a
}

func (a *Array) Err(err error) conflow.LogArray {
	a.arr = a.arr.Err(err)
	return a
}

func (a *Array) Bool(b bool) conflow.LogArray {
	a.arr = a.arr.Bool(b)
	return a
}

func (a *Array) Int(i int) conflow.LogArray {
	a.arr = a.arr.Int(i)
	return a
}

func (a *Array) Int8(i int8) conflow.LogArray {
	a.arr = a.arr.Int8(i)
	return a
}

func (a *Array) Int16(i int16) conflow.LogArray {
	a.arr = a.arr.Int16(i)
	return a
}

func (a *Array) Int32(i int32) conflow.LogArray {
	a.arr = a.arr.Int32(i)
	return a
}

func (a *Array) Int64(i int64) conflow.LogArray {
	a.arr = a.arr.Int64(i)
	return a
}

func (a *Array) Uint(i uint) conflow.LogArray {
	a.arr = a.arr.Uint(i)
	return a
}

func (a *Array) Uint8(i uint8) conflow.LogArray {
	a.arr = a.arr.Uint8(i)
	return a
}

func (a *Array) Uint16(i uint16) conflow.LogArray {
	a.arr = a.arr.Uint16(i)
	return a
}

func (a *Array) Uint32(i uint32) conflow.LogArray {
	a.arr = a.arr.Uint32(i)
	return a
}

func (a *Array) Uint64(i uint64) conflow.LogArray {
	a.arr = a.arr.Uint64(i)
	return a
}

func (a *Array) Float32(f float32) conflow.LogArray {
	a.arr = a.arr.Float32(f)
	return a
}

func (a *Array) Float64(f float64) conflow.LogArray {
	a.arr = a.arr.Float64(f)
	return a
}

func (a *Array) Time(t time.Time) conflow.LogArray {
	a.arr = a.arr.Time(t)
	return a
}

func (a *Array) Dur(d time.Duration) conflow.LogArray {
	a.arr = a.arr.Dur(d)
	return a
}

func (a *Array) Interface(i interface{}) conflow.LogArray {
	a.arr = a.arr.Interface(i)
	return a
}

func (a *Array) IPAddr(ip net.IP) conflow.LogArray {
	a.arr = a.arr.IPAddr(ip)
	return a
}

func (a *Array) IPPrefix(pfx net.IPNet) conflow.LogArray {
	a.arr = a.arr.IPPrefix(pfx)
	return a
}

func (a *Array) MACAddr(ha net.HardwareAddr) conflow.LogArray {
	a.arr = a.arr.MACAddr(ha)
	return a
}

type ArrayMarshalerWrapper struct {
	arr conflow.LogArrayMarshaler
}

func (a *ArrayMarshalerWrapper) MarshalZerologArray(arr *zerolog.Array) {
	a.arr.MarshalLogArray(&Array{arr: arr})
}
