package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

type ILogger interface {
	Debug(msg string, params ...any) string
	Info(msg string, params ...any) string
	Success(msg string, params ...any) string
	Warn(msg string, params ...any) string
	Error(msg string, params ...any) string
	Fatal(msg string, params ...any)
}

type logger struct {
	logStructureFormat string
	colorReset         string
	colorCyan          string
	colorRed           string
	colorYellow        string
	colorGreen         string
}

func NewLogger() ILogger {
	return &logger{
		logStructureFormat: "level:[%v] date:[%v] origin:[%s:%d] msg:[%v] ",
		colorReset:         "\033[0m",
		colorCyan:          "\033[96m",
		colorRed:           "\033[91m",
		colorYellow:        "\033[93m",
		colorGreen:         "\033[92m",
	}
}

func (l *logger) getCurrentTimeFormated() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l *logger) log(level, color, msg string, params ...any) string {
	pc, _, line, _ := runtime.Caller(2)
	logResult := fmt.Sprintf(l.logStructureFormat, level, l.getCurrentTimeFormated(), filepath.Base(runtime.FuncForPC(pc).Name()), line, msg)
	log.Printf(string(color)+logResult+string(l.colorReset), params...)
	return logResult
}

func (l *logger) Debug(msg string, params ...any) string {
	return l.log("DEBUG", string(l.colorReset), msg, params...)
}

func (l *logger) Info(msg string, params ...any) string {
	return l.log("INFO", string(l.colorCyan), msg, params...)
}

func (l *logger) Success(msg string, params ...any) string {
	return l.log("SUCCESS", string(l.colorGreen), msg, params...)
}

func (l *logger) Warn(msg string, params ...any) string {
	return l.log("WARN", string(l.colorYellow), msg, params...)
}

func (l *logger) Error(msg string, params ...any) string {
	return l.log("ERROR", string(l.colorRed), msg, params...)
}

func (l *logger) Fatal(msg string, params ...any) {
	logResult := l.log("FATAL", string(l.colorRed), msg, params...)
	log.Panic(logResult)
}
