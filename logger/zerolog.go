package logger

import (
	"github.com/opsidian/basil/basil"
	"github.com/rs/zerolog"
)

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
	return &ZeroLogEvent{event: z.logger.Debug()}
}

func (z *ZeroLogLogger) Info() basil.LogEvent {
	return &ZeroLogEvent{event: z.logger.Info()}
}

func (z *ZeroLogLogger) Warn() basil.LogEvent {
	return &ZeroLogEvent{event: z.logger.Warn()}
}

func (z *ZeroLogLogger) Error() basil.LogEvent {
	return &ZeroLogEvent{event: z.logger.Error()}
}

func (z *ZeroLogLogger) Fatal() basil.LogEvent {
	return &ZeroLogEvent{event: z.logger.Fatal()}
}

func (z *ZeroLogLogger) Panic() basil.LogEvent {
	return &ZeroLogEvent{event: z.logger.Panic()}
}

func (z *ZeroLogLogger) Log() basil.LogEvent {
	return &ZeroLogEvent{event: z.logger.Log()}
}

func (z *ZeroLogLogger) Print(v ...interface{}) {
	z.logger.Print(v...)
}

func (z *ZeroLogLogger) Printf(format string, v ...interface{}) {
	z.logger.Printf(format, v...)
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

type ZeroLogEvent struct {
	event *zerolog.Event
}

func (z *ZeroLogEvent) Fields(fields map[string]interface{}) basil.LogEvent {
	return &ZeroLogEvent{
		event: z.event.Fields(fields),
	}
}

func (z *ZeroLogEvent) Msg(msg string) {
	z.event.Msg(msg)
}

func (z *ZeroLogEvent) Msgf(format string, v ...interface{}) {
	z.event.Msgf(format, v...)
}
