package utils

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		info: log.New(os.Stdout, "koyo-site: [INFO] ", 0),
		warn: log.New(os.Stdout, "koyo-site: [WARN] ", 0),
		err:  log.New(os.Stderr, "koyo-site: [ERROR] ", 0),
	}
}

func (l *Logger) Info(format string, v ...any) {
	l.info.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...any) {
	l.warn.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...any) {
	l.err.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(format string, v ...any) {
	l.err.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
