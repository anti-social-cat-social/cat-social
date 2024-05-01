package logger

import (
	"fmt"
	"log/slog"
	"runtime"
)

func Info(message string) {
	slog.Info(message)
}

func Error(err error) {
	_, filename, line, _ := runtime.Caller(1)
	slog.Error(err.Error(), slog.String("file", fmt.Sprintf("%s:%d", filename, line)))
}
