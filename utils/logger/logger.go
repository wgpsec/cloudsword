package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Logger struct{}

var Println = &Logger{}
var Printf = &Logger{}

func (l *Logger) logWithNewline(level string, levelColor *color.Color, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelInfo := levelColor.Sprintf("[%s] %s", level, timestamp)
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", levelInfo, message)
}

func (l *Logger) logWithoutNewline(level string, levelColor *color.Color, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelInfo := levelColor.Sprintf("[%s] %s", level, timestamp)
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s", levelInfo, message)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.logWithNewline("DEBUG", color.New(color.FgBlue), format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.logWithNewline("INFO", color.New(color.FgGreen), format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.logWithNewline("WARN", color.New(color.FgYellow), format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.logWithNewline("ERROR", color.New(color.FgRed), format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logWithNewline("FATAL", color.New(color.FgMagenta), format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logWithoutNewline("DEBUG", color.New(color.FgBlue), format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logWithoutNewline("INFO", color.New(color.FgGreen), format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logWithoutNewline("WARN", color.New(color.FgYellow), format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logWithoutNewline("ERROR", color.New(color.FgRed), format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logWithoutNewline("FATAL", color.New(color.FgMagenta), format, args...)
}
