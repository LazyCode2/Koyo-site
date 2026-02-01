package utils

import (
	"fmt"
	"log"
	"os"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

type Logger struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		info: log.New(os.Stdout, bold+blue+"koyo-site: [INFO] "+reset, 0),
		warn: log.New(os.Stdout, bold+yellow+"koyo-site: [WARN] "+reset, 0),
		err:  log.New(os.Stderr, bold+red+"koyo-site: [ERROR] "+reset, 0),
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
