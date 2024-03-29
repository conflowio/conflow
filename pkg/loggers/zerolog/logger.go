// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package zerolog

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/conflowio/conflow/pkg/conflow"
)

func init() {
	zerolog.CallerSkipFrameCount = 3
	zerolog.TimeFieldFormat = conflow.LogTimeFormat
}

var nilEvent *Event

type Logger struct {
	logger zerolog.Logger
}

func NewLogger(logger zerolog.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func NewConsoleLogger(level zerolog.Level) *Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}).With().
		Timestamp().
		Logger().
		Level(level)
	return NewLogger(logger)
}

func NewDisabledLogger() *Logger {
	return NewLogger(zerolog.New(os.Stderr).Level(zerolog.Disabled))
}

func (l *Logger) With() conflow.LoggerContext {
	return &Context{
		context: l.logger.With(),
	}
}

func (l *Logger) Debug() conflow.LogEvent {
	if e := l.logger.Debug(); e != nil {
		return &Event{e: e}
	}
	return nilEvent
}

func (l *Logger) Info() conflow.LogEvent {
	if e := l.logger.Info(); e != nil {
		return &Event{e: e}
	}
	return nilEvent
}

func (l *Logger) Warn() conflow.LogEvent {
	if e := l.logger.Warn(); e != nil {
		return &Event{e: e}
	}
	return nilEvent
}

func (l *Logger) Error() conflow.LogEvent {
	if e := l.logger.Error(); e != nil {
		return &Event{e: e}
	}
	return nilEvent
}

func (l *Logger) Fatal() conflow.LogEvent {
	if e := l.logger.Fatal(); e != nil {
		return &Event{e: e}
	}
	return nilEvent
}

func (l *Logger) Panic() conflow.LogEvent {
	if e := l.logger.Panic(); e != nil {
		return &Event{e: e}
	}
	return nilEvent
}

func (l *Logger) Log() conflow.LogEvent {
	if e := l.logger.Log(); e != nil {
		return &Event{e: e}
	}
	return nilEvent
}

func (l *Logger) Print(v ...interface{}) {
	if e := l.Debug(); e.Enabled() {
		e.Msg(fmt.Sprint(v...))
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if e := l.Debug(); e.Enabled() {
		e.Msg(fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Array() conflow.LogArray {
	return &Array{
		arr: zerolog.Arr(),
	}
}
