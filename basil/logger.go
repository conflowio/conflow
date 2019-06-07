// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

type LogFields map[string]interface{}

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
}

// LoggerContext is an interface for setting up child loggers
type LoggerContext interface {
	Logger() Logger
	Fields(fields map[string]interface{}) LoggerContext
}

// LogEvent is an interface for enriching and sending log events
type LogEvent interface {
	Fields(fields map[string]interface{}) LogEvent
	Msg(msg string)
	Msgf(format string, v ...interface{})
}
