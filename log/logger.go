package log

import (
	"fmt"
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
	context map[string]any
}

func (l *Logger) chainLog(e *zerolog.Event, message string, tags ...any) {

	if tags != nil {
		logger := l.C(tags...)

		for key, value := range logger.context {
			e = e.Str(key, fmt.Sprintf("%v", value))
		}
	}

	e.Msg(message)
}

func (l *Logger) C(tags ...any) *Logger {

	logger := &Logger{
		Logger:  l.Logger,
		context: make(map[string]any),
	}

	for tag := range l.context {
		logger.context[tag] = l.context[tag]
	}

	if len(tags) == 0 {
		return logger
	}

	if len(tags)%2 == 1 {
		tags = append(tags, "<MISSINGARG>")
		l.chainLog(l.Error(), "logger.With: odd argument count", "tags", tags)
	}

	for i := 0; i < len(tags); i += 2 {
		tag, ok := tags[i].(string)
		if !ok {
			panic(fmt.Sprintf("logging tag is not a string. tag: %v", tag))
		}
		logger.context[tag] = tags[i+1]
	}

	return logger
}

func (l *Logger) With(tags ...any) {

	if len(tags)%2 == 1 {
		tags = append(tags, "<MISSINGARG>")
		l.chainLog(l.Error(), "logger.With: odd argument count", "tags", tags)
	}

	for i := 0; i < len(tags); i += 2 {
		tag, ok := tags[i].(string)
		if !ok {
			panic(fmt.Sprintf("logging tag is not a string. tag: %v", tag))
		}
		l.context[tag] = tags[i+1]
	}
}

func (l *Logger) D(message string, tags ...any) {
	l.chainLog(l.Debug(), message, tags...)
}

func (l *Logger) F(message string, tags ...any) {
	l.chainLog(l.Fatal(), message, tags...)
}

func (l *Logger) E(message string, tags ...any) {
	l.chainLog(l.Error(), message, tags...)
}

func (l *Logger) W(message string, tags ...any) {
	l.chainLog(l.Warn(), message, tags...)
}

func (l *Logger) L(message string, tags ...any) {
	l.chainLog(l.Log(), message, tags...)
}

func (l *Logger) I(message string, tags ...any) {
	l.chainLog(l.Info(), message, tags...)
}
