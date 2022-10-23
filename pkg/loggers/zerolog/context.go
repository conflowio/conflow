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

type Context struct {
	context zerolog.Context
}

func (c *Context) Logger() conflow.Logger {
	return &Logger{
		logger: c.context.Logger(),
	}
}

func (c *Context) Fields(fields map[string]interface{}) conflow.LoggerContext {
	c.context = c.context.Fields(fields)
	return c
}

func (c *Context) Dict(key string, dict conflow.LogEvent) conflow.LoggerContext {
	c.context = c.context.Dict(key, dict.(*Event).e)
	return c
}

func (c *Context) Array(key string, arr conflow.LogArrayMarshaler) conflow.LoggerContext {
	if a, ok := arr.(*Array); ok {
		c.context = c.context.Array(key, a.arr)
	} else {
		c.context = c.context.Array(key, &ArrayMarshalerWrapper{arr})
	}
	return c
}

func (c *Context) Object(key string, obj conflow.LogObjectMarshaler) conflow.LoggerContext {
	c.context = c.context.Object(key, &ObjectMarshalerWrapper{obj})
	return c
}

func (c *Context) EmbedObject(obj conflow.LogObjectMarshaler) conflow.LoggerContext {
	c.context = c.context.EmbedObject(&ObjectMarshalerWrapper{obj: obj})
	return c
}

func (c *Context) ID(key string, val conflow.ID) conflow.LoggerContext {
	c.context = c.context.Str(key, string(val))
	return c
}

func (c *Context) Str(key, val string) conflow.LoggerContext {
	c.context = c.context.Str(key, val)
	return c
}

func (c *Context) Strs(key string, vals []string) conflow.LoggerContext {
	c.context = c.context.Strs(key, vals)
	return c
}

func (c *Context) Bytes(key string, val []byte) conflow.LoggerContext {
	c.context = c.context.Bytes(key, val)
	return c
}

func (c *Context) Hex(key string, val []byte) conflow.LoggerContext {
	c.context = c.context.Hex(key, val)
	return c
}

func (c *Context) RawJSON(key string, b []byte) conflow.LoggerContext {
	c.context = c.context.RawJSON(key, b)
	return c
}

func (c *Context) AnErr(key string, err error) conflow.LoggerContext {
	c.context = c.context.AnErr(key, err)
	return c
}

func (c *Context) Errs(key string, errs []error) conflow.LoggerContext {
	c.context = c.context.Errs(key, errs)
	return c
}

func (c *Context) Err(err error) conflow.LoggerContext {
	c.context = c.context.Err(err)
	return c
}

func (c *Context) Bool(key string, b bool) conflow.LoggerContext {
	c.context = c.context.Bool(key, b)
	return c
}

func (c *Context) Bools(key string, b []bool) conflow.LoggerContext {
	c.context = c.context.Bools(key, b)
	return c
}

func (c *Context) Int(key string, i int) conflow.LoggerContext {
	c.context = c.context.Int(key, i)
	return c
}

func (c *Context) Ints(key string, i []int) conflow.LoggerContext {
	c.context = c.context.Ints(key, i)
	return c
}

func (c *Context) Int8(key string, i int8) conflow.LoggerContext {
	c.context = c.context.Int8(key, i)
	return c
}

func (c *Context) Ints8(key string, i []int8) conflow.LoggerContext {
	c.context = c.context.Ints8(key, i)
	return c
}

func (c *Context) Int16(key string, i int16) conflow.LoggerContext {
	c.context = c.context.Int16(key, i)
	return c
}

func (c *Context) Ints16(key string, i []int16) conflow.LoggerContext {
	c.context = c.context.Ints16(key, i)
	return c
}

func (c *Context) Int32(key string, i int32) conflow.LoggerContext {
	c.context = c.context.Int32(key, i)
	return c
}

func (c *Context) Ints32(key string, i []int32) conflow.LoggerContext {
	c.context = c.context.Ints32(key, i)
	return c
}

func (c *Context) Int64(key string, i int64) conflow.LoggerContext {
	c.context = c.context.Int64(key, i)
	return c
}

func (c *Context) Ints64(key string, i []int64) conflow.LoggerContext {
	c.context = c.context.Ints64(key, i)
	return c
}

func (c *Context) Uint(key string, i uint) conflow.LoggerContext {
	c.context = c.context.Uint(key, i)
	return c
}

func (c *Context) Uints(key string, i []uint) conflow.LoggerContext {
	c.context = c.context.Uints(key, i)
	return c
}

func (c *Context) Uint8(key string, i uint8) conflow.LoggerContext {
	c.context = c.context.Uint8(key, i)
	return c
}

func (c *Context) Uints8(key string, i []uint8) conflow.LoggerContext {
	c.context = c.context.Uints8(key, i)
	return c
}

func (c *Context) Uint16(key string, i uint16) conflow.LoggerContext {
	c.context = c.context.Uint16(key, i)
	return c
}

func (c *Context) Uints16(key string, i []uint16) conflow.LoggerContext {
	c.context = c.context.Uints16(key, i)
	return c
}

func (c *Context) Uint32(key string, i uint32) conflow.LoggerContext {
	c.context = c.context.Uint32(key, i)
	return c
}

func (c *Context) Uints32(key string, i []uint32) conflow.LoggerContext {
	c.context = c.context.Uints32(key, i)
	return c
}

func (c *Context) Uint64(key string, i uint64) conflow.LoggerContext {
	c.context = c.context.Uint64(key, i)
	return c
}

func (c *Context) Uints64(key string, i []uint64) conflow.LoggerContext {
	c.context = c.context.Uints64(key, i)
	return c
}

func (c *Context) Float32(key string, f float32) conflow.LoggerContext {
	c.context = c.context.Float32(key, f)
	return c
}

func (c *Context) Floats32(key string, f []float32) conflow.LoggerContext {
	c.context = c.context.Floats32(key, f)
	return c
}

func (c *Context) Float64(key string, f float64) conflow.LoggerContext {
	c.context = c.context.Float64(key, f)
	return c
}

func (c *Context) Floats64(key string, f []float64) conflow.LoggerContext {
	c.context = c.context.Floats64(key, f)
	return c
}

func (c *Context) Timestamp() conflow.LoggerContext {
	c.context = c.context.Timestamp()
	return c
}

func (c *Context) Time(key string, t time.Time) conflow.LoggerContext {
	c.context = c.context.Time(key, t)
	return c
}

func (c *Context) Times(key string, t []time.Time) conflow.LoggerContext {
	c.context = c.context.Times(key, t)
	return c
}

func (c *Context) Dur(key string, d time.Duration) conflow.LoggerContext {
	c.context = c.context.Dur(key, d)
	return c
}

func (c *Context) Durs(key string, d []time.Duration) conflow.LoggerContext {
	c.context = c.context.Durs(key, d)
	return c
}

func (c *Context) Interface(key string, i interface{}) conflow.LoggerContext {
	c.context = c.context.Interface(key, i)
	return c
}

func (c *Context) Caller() conflow.LoggerContext {
	c.context = c.context.Caller()
	return c
}

func (c *Context) CallerWithSkipFrameCount(skipFrameCount int) conflow.LoggerContext {
	c.context = c.context.CallerWithSkipFrameCount(skipFrameCount)
	return c
}

func (c *Context) Stack() conflow.LoggerContext {
	c.context = c.context.Stack()
	return c
}

func (c *Context) IPAddr(key string, ip net.IP) conflow.LoggerContext {
	c.context = c.context.IPAddr(key, ip)
	return c
}

func (c *Context) IPPrefix(key string, pfx net.IPNet) conflow.LoggerContext {
	c.context = c.context.IPPrefix(key, pfx)
	return c
}

func (c *Context) MACAddr(key string, ha net.HardwareAddr) conflow.LoggerContext {
	c.context = c.context.MACAddr(key, ha)
	return c
}
