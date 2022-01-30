// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"net"
	"time"
)

const LogTimeFormat = "2006-01-02T15:04:05.000Z07:00"

// Logger is an interface for structured logging
type Logger interface {
	With() LoggerContext
	Debug() LogEvent
	Info() LogEvent
	Warn() LogEvent
	Error() LogEvent
	Fatal() LogEvent
	Panic() LogEvent
	Log() LogEvent
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Array() LogArray
}

// LoggerContext is an interface for setting up child loggers
type LoggerContext interface {
	Logger() Logger
	Fields(fields map[string]interface{}) LoggerContext
	Dict(key string, dict LogEvent) LoggerContext
	Array(key string, arr LogArrayMarshaler) LoggerContext
	Object(key string, obj LogObjectMarshaler) LoggerContext
	EmbedObject(obj LogObjectMarshaler) LoggerContext
	ID(key string, val ID) LoggerContext
	Str(key, val string) LoggerContext
	Strs(key string, vals []string) LoggerContext
	Bytes(key string, val []byte) LoggerContext
	Hex(key string, val []byte) LoggerContext
	RawJSON(key string, b []byte) LoggerContext
	AnErr(key string, err error) LoggerContext
	Errs(key string, errs []error) LoggerContext
	Err(err error) LoggerContext
	Bool(key string, b bool) LoggerContext
	Bools(key string, b []bool) LoggerContext
	Int(key string, i int) LoggerContext
	Ints(key string, i []int) LoggerContext
	Int8(key string, i int8) LoggerContext
	Ints8(key string, i []int8) LoggerContext
	Int16(key string, i int16) LoggerContext
	Ints16(key string, i []int16) LoggerContext
	Int32(key string, i int32) LoggerContext
	Ints32(key string, i []int32) LoggerContext
	Int64(key string, i int64) LoggerContext
	Ints64(key string, i []int64) LoggerContext
	Uint(key string, i uint) LoggerContext
	Uints(key string, i []uint) LoggerContext
	Uint8(key string, i uint8) LoggerContext
	Uints8(key string, i []uint8) LoggerContext
	Uint16(key string, i uint16) LoggerContext
	Uints16(key string, i []uint16) LoggerContext
	Uint32(key string, i uint32) LoggerContext
	Uints32(key string, i []uint32) LoggerContext
	Uint64(key string, i uint64) LoggerContext
	Uints64(key string, i []uint64) LoggerContext
	Float32(key string, f float32) LoggerContext
	Floats32(key string, f []float32) LoggerContext
	Float64(key string, f float64) LoggerContext
	Floats64(key string, f []float64) LoggerContext
	Timestamp() LoggerContext
	Time(key string, t time.Time) LoggerContext
	Times(key string, t []time.Time) LoggerContext
	Dur(key string, d time.Duration) LoggerContext
	Durs(key string, d []time.Duration) LoggerContext
	Interface(key string, i interface{}) LoggerContext
	Caller() LoggerContext
	CallerWithSkipFrameCount(skipFrameCount int) LoggerContext
	Stack() LoggerContext
	IPAddr(key string, ip net.IP) LoggerContext
	IPPrefix(key string, pfx net.IPNet) LoggerContext
	MACAddr(key string, ha net.HardwareAddr) LoggerContext
}

// LogEvent is an interface for enriching and sending log events
type LogEvent interface {
	Enabled() bool
	Discard() LogEvent
	Msg(msg string)
	Msgf(format string, v ...interface{})
	Fields(fields map[string]interface{}) LogEvent
	Dict(key string, dict LogEvent) LogEvent
	Array(key string, arr LogArrayMarshaler) LogEvent
	Object(key string, obj LogObjectMarshaler) LogEvent
	EmbedObject(obj LogObjectMarshaler) LogEvent
	ID(key string, val ID) LogEvent
	Str(key, val string) LogEvent
	Strs(key string, vals []string) LogEvent
	Bytes(key string, val []byte) LogEvent
	Hex(key string, val []byte) LogEvent
	RawJSON(key string, b []byte) LogEvent
	AnErr(key string, err error) LogEvent
	Errs(key string, errs []error) LogEvent
	Err(err error) LogEvent
	Stack() LogEvent
	Bool(key string, b bool) LogEvent
	Bools(key string, b []bool) LogEvent
	Int(key string, i int) LogEvent
	Ints(key string, i []int) LogEvent
	Int8(key string, i int8) LogEvent
	Ints8(key string, i []int8) LogEvent
	Int16(key string, i int16) LogEvent
	Ints16(key string, i []int16) LogEvent
	Int32(key string, i int32) LogEvent
	Ints32(key string, i []int32) LogEvent
	Int64(key string, i int64) LogEvent
	Ints64(key string, i []int64) LogEvent
	Uint(key string, i uint) LogEvent
	Uints(key string, i []uint) LogEvent
	Uint8(key string, i uint8) LogEvent
	Uints8(key string, i []uint8) LogEvent
	Uint16(key string, i uint16) LogEvent
	Uints16(key string, i []uint16) LogEvent
	Uint32(key string, i uint32) LogEvent
	Uints32(key string, i []uint32) LogEvent
	Uint64(key string, i uint64) LogEvent
	Uints64(key string, i []uint64) LogEvent
	Float32(key string, f float32) LogEvent
	Floats32(key string, f []float32) LogEvent
	Float64(key string, f float64) LogEvent
	Floats64(key string, f []float64) LogEvent
	Timestamp() LogEvent
	Time(key string, t time.Time) LogEvent
	Times(key string, t []time.Time) LogEvent
	Dur(key string, d time.Duration) LogEvent
	Durs(key string, d []time.Duration) LogEvent
	TimeDiff(key string, t time.Time, start time.Time) LogEvent
	Interface(key string, i interface{}) LogEvent
	Caller() LogEvent
	IPAddr(key string, ip net.IP) LogEvent
	IPPrefix(key string, pfx net.IPNet) LogEvent
	MACAddr(key string, ha net.HardwareAddr) LogEvent
}

type LogObjectMarshaler interface {
	MarshalLogObject(e LogEvent)
}

type LogArrayMarshaler interface {
	MarshalLogArray(a LogArray)
}
type LogArray interface {
	LogArrayMarshaler
	Object(obj LogObjectMarshaler) LogArray
	ID(val ID) LogArray
	Str(val string) LogArray
	Bytes(val []byte) LogArray
	Hex(val []byte) LogArray
	Err(err error) LogArray
	Bool(b bool) LogArray
	Int(i int) LogArray
	Int8(i int8) LogArray
	Int16(i int16) LogArray
	Int32(i int32) LogArray
	Int64(i int64) LogArray
	Uint(i uint) LogArray
	Uint8(i uint8) LogArray
	Uint16(i uint16) LogArray
	Uint32(i uint32) LogArray
	Uint64(i uint64) LogArray
	Float32(f float32) LogArray
	Float64(f float64) LogArray
	Time(t time.Time) LogArray
	Dur(d time.Duration) LogArray
	Interface(i interface{}) LogArray
	IPAddr(ip net.IP) LogArray
	IPPrefix(pfx net.IPNet) LogArray
	MACAddr(ha net.HardwareAddr) LogArray
}
