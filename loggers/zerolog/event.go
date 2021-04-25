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

	"github.com/opsidian/basil/basil"
)

type Event struct {
	e *zerolog.Event
}

func (e *Event) Enabled() bool {
	return e != nil && e.e.Enabled()
}

func (e *Event) Discard() basil.LogEvent {
	if e == nil {
		return e
	}
	e.e.Discard()
	return nil
}

func (e *Event) Msg(msg string) {
	if e == nil {
		return
	}
	e.e.Msg(msg)
}

func (e *Event) Msgf(format string, v ...interface{}) {
	if e == nil {
		return
	}
	e.e.Msgf(format, v...)
}

func (e *Event) Fields(fields map[string]interface{}) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Fields(fields)
	return e
}

func (e *Event) Dict(key string, dict basil.LogEvent) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Dict(key, dict.(*Event).e)
	return e
}

func (e *Event) Array(key string, arr basil.LogArrayMarshaler) basil.LogEvent {
	if a, ok := arr.(*Array); ok {
		e.e = e.e.Array(key, a.arr)
	} else {
		e.e = e.e.Array(key, &ArrayMarshalerWrapper{arr})
	}
	return e
}

func (e *Event) Object(key string, obj basil.LogObjectMarshaler) basil.LogEvent {
	e.e = e.e.Object(key, &ObjectMarshalerWrapper{obj})
	return e
}

func (e *Event) EmbedObject(obj basil.LogObjectMarshaler) basil.LogEvent {
	if e == nil {
		return e
	}
	obj.MarshalLogObject(e)
	return e
}

func (e *Event) ID(key string, val basil.ID) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Str(key, string(val))
	return e
}

func (e *Event) Str(key, val string) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Str(key, val)
	return e
}

func (e *Event) Strs(key string, vals []string) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Strs(key, vals)
	return e
}

func (e *Event) Bytes(key string, val []byte) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Bytes(key, val)
	return e
}

func (e *Event) Hex(key string, val []byte) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Hex(key, val)
	return e
}

func (e *Event) RawJSON(key string, b []byte) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.RawJSON(key, b)
	return e
}

func (e *Event) AnErr(key string, err error) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.AnErr(key, err)
	return e
}

func (e *Event) Errs(key string, errs []error) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Errs(key, errs)
	return e
}

func (e *Event) Err(err error) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Err(err)
	return e
}

func (e *Event) Stack() basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Stack()
	return e
}

func (e *Event) Bool(key string, b bool) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Bool(key, b)
	return e
}

func (e *Event) Bools(key string, b []bool) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Bools(key, b)
	return e
}

func (e *Event) Int(key string, i int) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Int(key, i)
	return e
}

func (e *Event) Ints(key string, i []int) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Ints(key, i)
	return e
}

func (e *Event) Int8(key string, i int8) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Int8(key, i)
	return e
}

func (e *Event) Ints8(key string, i []int8) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Ints8(key, i)
	return e
}

func (e *Event) Int16(key string, i int16) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Int16(key, i)
	return e
}

func (e *Event) Ints16(key string, i []int16) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Ints16(key, i)
	return e
}

func (e *Event) Int32(key string, i int32) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Int32(key, i)
	return e
}

func (e *Event) Ints32(key string, i []int32) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Ints32(key, i)
	return e
}

func (e *Event) Int64(key string, i int64) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Int64(key, i)
	return e
}

func (e *Event) Ints64(key string, i []int64) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Ints64(key, i)
	return e
}

func (e *Event) Uint(key string, i uint) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uint(key, i)
	return e
}

func (e *Event) Uints(key string, i []uint) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uints(key, i)
	return e
}

func (e *Event) Uint8(key string, i uint8) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uint8(key, i)
	return e
}

func (e *Event) Uints8(key string, i []uint8) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uints8(key, i)
	return e
}

func (e *Event) Uint16(key string, i uint16) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uint16(key, i)
	return e
}

func (e *Event) Uints16(key string, i []uint16) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uints16(key, i)
	return e
}

func (e *Event) Uint32(key string, i uint32) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uint32(key, i)
	return e
}

func (e *Event) Uints32(key string, i []uint32) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uints32(key, i)
	return e
}

func (e *Event) Uint64(key string, i uint64) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uint64(key, i)
	return e
}

func (e *Event) Uints64(key string, i []uint64) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Uints64(key, i)
	return e
}

func (e *Event) Float32(key string, f float32) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Float32(key, f)
	return e
}

func (e *Event) Floats32(key string, f []float32) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Floats32(key, f)
	return e
}

func (e *Event) Float64(key string, f float64) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Float64(key, f)
	return e
}

func (e *Event) Floats64(key string, f []float64) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Floats64(key, f)
	return e
}

func (e *Event) Timestamp() basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Timestamp()
	return e
}

func (e *Event) Time(key string, t time.Time) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Time(key, t)
	return e
}

func (e *Event) Times(key string, t []time.Time) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Times(key, t)
	return e
}

func (e *Event) Dur(key string, d time.Duration) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Dur(key, d)
	return e
}

func (e *Event) Durs(key string, d []time.Duration) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Durs(key, d)
	return e
}

func (e *Event) TimeDiff(key string, t time.Time, start time.Time) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.TimeDiff(key, t, start)
	return e
}

func (e *Event) Interface(key string, i interface{}) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Interface(key, i)
	return e
}

func (e *Event) Caller() basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.Caller()
	return e
}

func (e *Event) IPAddr(key string, ip net.IP) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.IPAddr(key, ip)
	return e
}

func (e *Event) IPPrefix(key string, pfx net.IPNet) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.IPPrefix(key, pfx)
	return e
}

func (e *Event) MACAddr(key string, ha net.HardwareAddr) basil.LogEvent {
	if e == nil {
		return e
	}
	e.e = e.e.MACAddr(key, ha)
	return e
}
