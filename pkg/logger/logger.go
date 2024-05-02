package logger

import (
	"fmt"
	"log/slog"
	"runtime"
)

func Info(message interface{}) {
	slog.Info(fmt.Sprintf("%v", message))
}

func Error(err error) {
	_, filename, line, _ := runtime.Caller(1)
	slog.Error(err.Error(), slog.String("file", fmt.Sprintf("%s:%d", filename, line)))
}
