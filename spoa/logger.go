package main

import "log"

type Logger struct {
	*log.Logger
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf(format, v)
}

func NewLogger(l *log.Logger) *Logger {
	return &Logger{l}
}
