// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package logger

import (
	"fmt"
	"net"
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.CallerSkipFrameCount = 3
	zerolog.TimeFieldFormat = basil.LogTimeFormat
}

var nilZerologEvent *ZeroLogEvent

type ZeroLogLogger struct {
	logger zerolog.Logger
}

func NewZeroLogLogger(logger zerolog.Logger) *ZeroLogLogger {
	return &ZeroLogLogger{
		logger: logger,
	}
}

func (z *ZeroLogLogger) With() basil.LoggerContext {
	return &ZeroLogContext{
		context: z.logger.With(),
	}
}

func (z *ZeroLogLogger) Debug() basil.LogEvent {
	if e := z.logger.Debug(); e != nil {
		return &ZeroLogEvent{e: e}
	}
	return nilZerologEvent
}

func (z *ZeroLogLogger) Info() basil.LogEvent {
	if e := z.logger.Info(); e != nil {
		return &ZeroLogEvent{e: e}
	}
	return nilZerologEvent
}

func (z *ZeroLogLogger) Warn() basil.LogEvent {
	if e := z.logger.Warn(); e != nil {
		return &ZeroLogEvent{e: e}
	}
	return nilZerologEvent
}

func (z *ZeroLogLogger) Error() basil.LogEvent {
	if e := z.logger.Error(); e != nil {
		return &ZeroLogEvent{e: e}
	}
	return nilZerologEvent
}

func (z *ZeroLogLogger) Fatal() basil.LogEvent {
	if e := z.logger.Fatal(); e != nil {
		return &ZeroLogEvent{e: e}
	}
	return nilZerologEvent
}

func (z *ZeroLogLogger) Panic() basil.LogEvent {
	if e := z.logger.Panic(); e != nil {
		return &ZeroLogEvent{e: e}
	}
	return nilZerologEvent
}

func (z *ZeroLogLogger) Log() basil.LogEvent {
	if e := z.logger.Log(); e != nil {
		return &ZeroLogEvent{e: e}
	}
	return nilZerologEvent
}

func (z *ZeroLogLogger) Print(v ...interface{}) {
	if e := z.Debug(); e.Enabled() {
		e.Msg(fmt.Sprint(v...))
	}
}

func (z *ZeroLogLogger) Printf(format string, v ...interface{}) {
	if e := z.Debug(); e.Enabled() {
		e.Msg(fmt.Sprintf(format, v...))
	}
}

func (z *ZeroLogLogger) Array() basil.LogArray {
	return &ZerologArray{
		arr: zerolog.Arr(),
	}
}

type ZeroLogContext struct {
	context zerolog.Context
}

func (z *ZeroLogContext) Logger() basil.Logger {
	return &ZeroLogLogger{
		logger: z.context.Logger(),
	}
}

func (z *ZeroLogContext) Fields(fields map[string]interface{}) basil.LoggerContext {
	z.context = z.context.Fields(fields)
	return z
}

func (z *ZeroLogContext) Dict(key string, dict basil.LogEvent) basil.LoggerContext {
	z.context = z.context.Dict(key, dict.(*ZeroLogEvent).e)
	return z
}

func (z *ZeroLogContext) Array(key string, arr basil.LogArrayMarshaler) basil.LoggerContext {
	if a, ok := arr.(*ZerologArray); ok {
		z.context = z.context.Array(key, a.arr)
	} else {
		z.context = z.context.Array(key, &ZeroLogArrayMarshalerWrapper{arr})
	}
	return z
}

func (z *ZeroLogContext) Object(key string, obj basil.LogObjectMarshaler) basil.LoggerContext {
	z.context = z.context.Object(key, &ZeroLogObjectMarshalerWrapper{obj})
	return z
}

func (z *ZeroLogContext) EmbedObject(obj basil.LogObjectMarshaler) basil.LoggerContext {
	z.context = z.context.EmbedObject(&ZeroLogObjectMarshalerWrapper{obj: obj})
	return z
}

func (z *ZeroLogContext) ID(key string, val basil.ID) basil.LoggerContext {
	z.context = z.context.Str(key, string(val))
	return z
}

func (z *ZeroLogContext) Str(key, val string) basil.LoggerContext {
	z.context = z.context.Str(key, val)
	return z
}

func (z *ZeroLogContext) Strs(key string, vals []string) basil.LoggerContext {
	z.context = z.context.Strs(key, vals)
	return z
}

func (z *ZeroLogContext) Bytes(key string, val []byte) basil.LoggerContext {
	z.context = z.context.Bytes(key, val)
	return z
}

func (z *ZeroLogContext) Hex(key string, val []byte) basil.LoggerContext {
	z.context = z.context.Hex(key, val)
	return z
}

func (z *ZeroLogContext) RawJSON(key string, b []byte) basil.LoggerContext {
	z.context = z.context.RawJSON(key, b)
	return z
}

func (z *ZeroLogContext) AnErr(key string, err error) basil.LoggerContext {
	z.context = z.context.AnErr(key, err)
	return z
}

func (z *ZeroLogContext) Errs(key string, errs []error) basil.LoggerContext {
	z.context = z.context.Errs(key, errs)
	return z
}

func (z *ZeroLogContext) Err(err error) basil.LoggerContext {
	z.context = z.context.Err(err)
	return z
}

func (z *ZeroLogContext) Bool(key string, b bool) basil.LoggerContext {
	z.context = z.context.Bool(key, b)
	return z
}

func (z *ZeroLogContext) Bools(key string, b []bool) basil.LoggerContext {
	z.context = z.context.Bools(key, b)
	return z
}

func (z *ZeroLogContext) Int(key string, i int) basil.LoggerContext {
	z.context = z.context.Int(key, i)
	return z
}

func (z *ZeroLogContext) Ints(key string, i []int) basil.LoggerContext {
	z.context = z.context.Ints(key, i)
	return z
}

func (z *ZeroLogContext) Int8(key string, i int8) basil.LoggerContext {
	z.context = z.context.Int8(key, i)
	return z
}

func (z *ZeroLogContext) Ints8(key string, i []int8) basil.LoggerContext {
	z.context = z.context.Ints8(key, i)
	return z
}

func (z *ZeroLogContext) Int16(key string, i int16) basil.LoggerContext {
	z.context = z.context.Int16(key, i)
	return z
}

func (z *ZeroLogContext) Ints16(key string, i []int16) basil.LoggerContext {
	z.context = z.context.Ints16(key, i)
	return z
}

func (z *ZeroLogContext) Int32(key string, i int32) basil.LoggerContext {
	z.context = z.context.Int32(key, i)
	return z
}

func (z *ZeroLogContext) Ints32(key string, i []int32) basil.LoggerContext {
	z.context = z.context.Ints32(key, i)
	return z
}

func (z *ZeroLogContext) Int64(key string, i int64) basil.LoggerContext {
	z.context = z.context.Int64(key, i)
	return z
}

func (z *ZeroLogContext) Ints64(key string, i []int64) basil.LoggerContext {
	z.context = z.context.Ints64(key, i)
	return z
}

func (z *ZeroLogContext) Uint(key string, i uint) basil.LoggerContext {
	z.context = z.context.Uint(key, i)
	return z
}

func (z *ZeroLogContext) Uints(key string, i []uint) basil.LoggerContext {
	z.context = z.context.Uints(key, i)
	return z
}

func (z *ZeroLogContext) Uint8(key string, i uint8) basil.LoggerContext {
	z.context = z.context.Uint8(key, i)
	return z
}

func (z *ZeroLogContext) Uints8(key string, i []uint8) basil.LoggerContext {
	z.context = z.context.Uints8(key, i)
	return z
}

func (z *ZeroLogContext) Uint16(key string, i uint16) basil.LoggerContext {
	z.context = z.context.Uint16(key, i)
	return z
}

func (z *ZeroLogContext) Uints16(key string, i []uint16) basil.LoggerContext {
	z.context = z.context.Uints16(key, i)
	return z
}

func (z *ZeroLogContext) Uint32(key string, i uint32) basil.LoggerContext {
	z.context = z.context.Uint32(key, i)
	return z
}

func (z *ZeroLogContext) Uints32(key string, i []uint32) basil.LoggerContext {
	z.context = z.context.Uints32(key, i)
	return z
}

func (z *ZeroLogContext) Uint64(key string, i uint64) basil.LoggerContext {
	z.context = z.context.Uint64(key, i)
	return z
}

func (z *ZeroLogContext) Uints64(key string, i []uint64) basil.LoggerContext {
	z.context = z.context.Uints64(key, i)
	return z
}

func (z *ZeroLogContext) Float32(key string, f float32) basil.LoggerContext {
	z.context = z.context.Float32(key, f)
	return z
}

func (z *ZeroLogContext) Floats32(key string, f []float32) basil.LoggerContext {
	z.context = z.context.Floats32(key, f)
	return z
}

func (z *ZeroLogContext) Float64(key string, f float64) basil.LoggerContext {
	z.context = z.context.Float64(key, f)
	return z
}

func (z *ZeroLogContext) Floats64(key string, f []float64) basil.LoggerContext {
	z.context = z.context.Floats64(key, f)
	return z
}

func (z *ZeroLogContext) Timestamp() basil.LoggerContext {
	z.context = z.context.Timestamp()
	return z
}

func (z *ZeroLogContext) Time(key string, t time.Time) basil.LoggerContext {
	z.context = z.context.Time(key, t)
	return z
}

func (z *ZeroLogContext) Times(key string, t []time.Time) basil.LoggerContext {
	z.context = z.context.Times(key, t)
	return z
}

func (z *ZeroLogContext) Dur(key string, d time.Duration) basil.LoggerContext {
	z.context = z.context.Dur(key, d)
	return z
}

func (z *ZeroLogContext) Durs(key string, d []time.Duration) basil.LoggerContext {
	z.context = z.context.Durs(key, d)
	return z
}

func (z *ZeroLogContext) Interface(key string, i interface{}) basil.LoggerContext {
	z.context = z.context.Interface(key, i)
	return z
}

func (z *ZeroLogContext) Caller() basil.LoggerContext {
	z.context = z.context.Caller()
	return z
}

func (z *ZeroLogContext) CallerWithSkipFrameCount(skipFrameCount int) basil.LoggerContext {
	z.context = z.context.CallerWithSkipFrameCount(skipFrameCount)
	return z
}

func (z *ZeroLogContext) Stack() basil.LoggerContext {
	z.context = z.context.Stack()
	return z
}

func (z *ZeroLogContext) IPAddr(key string, ip net.IP) basil.LoggerContext {
	z.context = z.context.IPAddr(key, ip)
	return z
}

func (z *ZeroLogContext) IPPrefix(key string, pfx net.IPNet) basil.LoggerContext {
	z.context = z.context.IPPrefix(key, pfx)
	return z
}

func (z *ZeroLogContext) MACAddr(key string, ha net.HardwareAddr) basil.LoggerContext {
	z.context = z.context.MACAddr(key, ha)
	return z
}

type ZeroLogEvent struct {
	e *zerolog.Event
}

func (z *ZeroLogEvent) Enabled() bool {
	return z != nil && z.e.Enabled()
}

func (z *ZeroLogEvent) Discard() basil.LogEvent {
	if z == nil {
		return z
	}
	z.e.Discard()
	return nil
}

func (z *ZeroLogEvent) Msg(msg string) {
	if z == nil {
		return
	}
	z.e.Msg(msg)
}

func (z *ZeroLogEvent) Msgf(format string, v ...interface{}) {
	if z == nil {
		return
	}
	z.e.Msgf(format, v...)
}

func (z *ZeroLogEvent) Fields(fields map[string]interface{}) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Fields(fields)
	return z
}

func (z *ZeroLogEvent) Dict(key string, dict basil.LogEvent) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Dict(key, dict.(*ZeroLogEvent).e)
	return z
}

func (z *ZeroLogEvent) Array(key string, arr basil.LogArrayMarshaler) basil.LogEvent {
	if a, ok := arr.(*ZerologArray); ok {
		z.e = z.e.Array(key, a.arr)
	} else {
		z.e = z.e.Array(key, &ZeroLogArrayMarshalerWrapper{arr})
	}
	return z
}

func (z *ZeroLogEvent) Object(key string, obj basil.LogObjectMarshaler) basil.LogEvent {
	z.e = z.e.Object(key, &ZeroLogObjectMarshalerWrapper{obj})
	return z
}

func (z *ZeroLogEvent) EmbedObject(obj basil.LogObjectMarshaler) basil.LogEvent {
	if z == nil {
		return z
	}
	obj.MarshalLogObject(z)
	return z
}

func (z *ZeroLogEvent) ID(key string, val basil.ID) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Str(key, string(val))
	return z
}

func (z *ZeroLogEvent) Str(key, val string) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Str(key, val)
	return z
}

func (z *ZeroLogEvent) Strs(key string, vals []string) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Strs(key, vals)
	return z
}

func (z *ZeroLogEvent) Bytes(key string, val []byte) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Bytes(key, val)
	return z
}

func (z *ZeroLogEvent) Hex(key string, val []byte) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Hex(key, val)
	return z
}

func (z *ZeroLogEvent) RawJSON(key string, b []byte) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.RawJSON(key, b)
	return z
}

func (z *ZeroLogEvent) AnErr(key string, err error) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.AnErr(key, err)
	return z
}

func (z *ZeroLogEvent) Errs(key string, errs []error) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Errs(key, errs)
	return z
}

func (z *ZeroLogEvent) Err(err error) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Err(err)
	return z
}

func (z *ZeroLogEvent) Stack() basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Stack()
	return z
}

func (z *ZeroLogEvent) Bool(key string, b bool) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Bool(key, b)
	return z
}

func (z *ZeroLogEvent) Bools(key string, b []bool) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Bools(key, b)
	return z
}

func (z *ZeroLogEvent) Int(key string, i int) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Int(key, i)
	return z
}

func (z *ZeroLogEvent) Ints(key string, i []int) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Ints(key, i)
	return z
}

func (z *ZeroLogEvent) Int8(key string, i int8) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Int8(key, i)
	return z
}

func (z *ZeroLogEvent) Ints8(key string, i []int8) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Ints8(key, i)
	return z
}

func (z *ZeroLogEvent) Int16(key string, i int16) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Int16(key, i)
	return z
}

func (z *ZeroLogEvent) Ints16(key string, i []int16) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Ints16(key, i)
	return z
}

func (z *ZeroLogEvent) Int32(key string, i int32) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Int32(key, i)
	return z
}

func (z *ZeroLogEvent) Ints32(key string, i []int32) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Ints32(key, i)
	return z
}

func (z *ZeroLogEvent) Int64(key string, i int64) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Int64(key, i)
	return z
}

func (z *ZeroLogEvent) Ints64(key string, i []int64) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Ints64(key, i)
	return z
}

func (z *ZeroLogEvent) Uint(key string, i uint) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uint(key, i)
	return z
}

func (z *ZeroLogEvent) Uints(key string, i []uint) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uints(key, i)
	return z
}

func (z *ZeroLogEvent) Uint8(key string, i uint8) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uint8(key, i)
	return z
}

func (z *ZeroLogEvent) Uints8(key string, i []uint8) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uints8(key, i)
	return z
}

func (z *ZeroLogEvent) Uint16(key string, i uint16) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uint16(key, i)
	return z
}

func (z *ZeroLogEvent) Uints16(key string, i []uint16) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uints16(key, i)
	return z
}

func (z *ZeroLogEvent) Uint32(key string, i uint32) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uint32(key, i)
	return z
}

func (z *ZeroLogEvent) Uints32(key string, i []uint32) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uints32(key, i)
	return z
}

func (z *ZeroLogEvent) Uint64(key string, i uint64) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uint64(key, i)
	return z
}

func (z *ZeroLogEvent) Uints64(key string, i []uint64) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Uints64(key, i)
	return z
}

func (z *ZeroLogEvent) Float32(key string, f float32) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Float32(key, f)
	return z
}

func (z *ZeroLogEvent) Floats32(key string, f []float32) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Floats32(key, f)
	return z
}

func (z *ZeroLogEvent) Float64(key string, f float64) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Float64(key, f)
	return z
}

func (z *ZeroLogEvent) Floats64(key string, f []float64) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Floats64(key, f)
	return z
}

func (z *ZeroLogEvent) Timestamp() basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Timestamp()
	return z
}

func (z *ZeroLogEvent) Time(key string, t time.Time) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Time(key, t)
	return z
}

func (z *ZeroLogEvent) Times(key string, t []time.Time) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Times(key, t)
	return z
}

func (z *ZeroLogEvent) Dur(key string, d time.Duration) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Dur(key, d)
	return z
}

func (z *ZeroLogEvent) Durs(key string, d []time.Duration) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Durs(key, d)
	return z
}

func (z *ZeroLogEvent) TimeDiff(key string, t time.Time, start time.Time) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.TimeDiff(key, t, start)
	return z
}

func (z *ZeroLogEvent) Interface(key string, i interface{}) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Interface(key, i)
	return z
}

func (z *ZeroLogEvent) Caller() basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.Caller()
	return z
}

func (z *ZeroLogEvent) IPAddr(key string, ip net.IP) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.IPAddr(key, ip)
	return z
}

func (z *ZeroLogEvent) IPPrefix(key string, pfx net.IPNet) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.IPPrefix(key, pfx)
	return z
}

func (z *ZeroLogEvent) MACAddr(key string, ha net.HardwareAddr) basil.LogEvent {
	if z == nil {
		return z
	}
	z.e = z.e.MACAddr(key, ha)
	return z
}

type ZeroLogObjectMarshalerWrapper struct {
	obj basil.LogObjectMarshaler
}

func (z *ZeroLogObjectMarshalerWrapper) MarshalZerologObject(e *zerolog.Event) {
	z.obj.MarshalLogObject(&ZeroLogEvent{e: e})
}

type ZerologArray struct {
	arr *zerolog.Array
}

func (z *ZerologArray) MarshalZerologArray(*zerolog.Array) {
}

func (z *ZerologArray) MarshalLogArray(basil.LogArray) {
}

func (z *ZerologArray) Object(obj basil.LogObjectMarshaler) basil.LogArray {
	z.arr = z.arr.Object(&ZeroLogObjectMarshalerWrapper{obj: obj})
	return z
}

func (z *ZerologArray) ID(val basil.ID) basil.LogArray {
	z.arr = z.arr.Str(string(val))
	return z
}

func (z *ZerologArray) Str(val string) basil.LogArray {
	z.arr = z.arr.Str(val)
	return z
}

func (z *ZerologArray) Bytes(val []byte) basil.LogArray {
	z.arr = z.arr.Bytes(val)
	return z
}

func (z *ZerologArray) Hex(val []byte) basil.LogArray {
	z.arr = z.arr.Hex(val)
	return z
}

func (z *ZerologArray) Err(err error) basil.LogArray {
	z.arr = z.arr.Err(err)
	return z
}

func (z *ZerologArray) Bool(b bool) basil.LogArray {
	z.arr = z.arr.Bool(b)
	return z
}

func (z *ZerologArray) Int(i int) basil.LogArray {
	z.arr = z.arr.Int(i)
	return z
}

func (z *ZerologArray) Int8(i int8) basil.LogArray {
	z.arr = z.arr.Int8(i)
	return z
}

func (z *ZerologArray) Int16(i int16) basil.LogArray {
	z.arr = z.arr.Int16(i)
	return z
}

func (z *ZerologArray) Int32(i int32) basil.LogArray {
	z.arr = z.arr.Int32(i)
	return z
}

func (z *ZerologArray) Int64(i int64) basil.LogArray {
	z.arr = z.arr.Int64(i)
	return z
}

func (z *ZerologArray) Uint(i uint) basil.LogArray {
	z.arr = z.arr.Uint(i)
	return z
}

func (z *ZerologArray) Uint8(i uint8) basil.LogArray {
	z.arr = z.arr.Uint8(i)
	return z
}

func (z *ZerologArray) Uint16(i uint16) basil.LogArray {
	z.arr = z.arr.Uint16(i)
	return z
}

func (z *ZerologArray) Uint32(i uint32) basil.LogArray {
	z.arr = z.arr.Uint32(i)
	return z
}

func (z *ZerologArray) Uint64(i uint64) basil.LogArray {
	z.arr = z.arr.Uint64(i)
	return z
}

func (z *ZerologArray) Float32(f float32) basil.LogArray {
	z.arr = z.arr.Float32(f)
	return z
}

func (z *ZerologArray) Float64(f float64) basil.LogArray {
	z.arr = z.arr.Float64(f)
	return z
}

func (z *ZerologArray) Time(t time.Time) basil.LogArray {
	z.arr = z.arr.Time(t)
	return z
}

func (z *ZerologArray) Dur(d time.Duration) basil.LogArray {
	z.arr = z.arr.Dur(d)
	return z
}

func (z *ZerologArray) Interface(i interface{}) basil.LogArray {
	z.arr = z.arr.Interface(i)
	return z
}

func (z *ZerologArray) IPAddr(ip net.IP) basil.LogArray {
	z.arr = z.arr.IPAddr(ip)
	return z
}

func (z *ZerologArray) IPPrefix(pfx net.IPNet) basil.LogArray {
	z.arr = z.arr.IPPrefix(pfx)
	return z
}

func (z *ZerologArray) MACAddr(ha net.HardwareAddr) basil.LogArray {
	z.arr = z.arr.MACAddr(ha)
	return z
}

type ZeroLogArrayMarshalerWrapper struct {
	arr basil.LogArrayMarshaler
}

func (z *ZeroLogArrayMarshalerWrapper) MarshalZerologArray(arr *zerolog.Array) {
	z.arr.MarshalLogArray(&ZerologArray{arr: arr})
}
