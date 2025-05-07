package flog

import (
	"fmt"
	"io"
)

func Info(format string, args ...interface{}) {
	if logger != nil {
		logger.log(LogInfo, format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if logger != nil {
		logger.log(LogError, format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if logger != nil {
		logger.log(LogWarn, format, args...)
	}
}

func Debug(format string, args ...interface{}) {
	if logger != nil {
		logger.log(LogDebug, format, args...)
	}
}

func Success(format string, args ...interface{}) {
	if logger != nil {
		logger.log(LogSuccess, format, args...)
	}
}

func Panic(format string, args ...interface{}) {
	if logger != nil {
		message := fmt.Sprintf(format, args...)
		msg := logger.formatLogEntry(LogPanic, message, nil)
		panic(msg)
	}
}

func Reader(format string, args ...interface{}) {
	if logger != nil {
		logger.log(LogReader, format, args...)
	}
}

func Channel(format string, args ...interface{}) {
	if logger != nil {
		logger.log(LogChannel, format, args...)
	}
}

func WatchReader(r io.Reader, level LogLevel) {
	if logger != nil {
		logger.WatchReader(r, level)
	}
}

func WatchChannel(ch interface{}, level LogLevel) {
	if logger != nil {
		logger.WatchChannel(ch, level)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogInfo, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogError, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogWarn, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogDebug, format, args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	l.log(LogSuccess, format, args...)
}

func (l *Logger) Panic(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	msg := l.formatLogEntry(LogPanic, message, nil)
	panic(msg)
}

func (l *Logger) Reader(format string, args ...interface{}) {
	l.log(LogReader, format, args...)
}

func (l *Logger) Channel(format string, args ...interface{}) {
	l.log(LogChannel, format, args...)
}
