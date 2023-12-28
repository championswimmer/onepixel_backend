package utils

import (
	"log"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	RedBold     = "\033[31;1m"
	GreenBold   = "\033[32;1m"
	YellowBold  = "\033[33;1m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	CyanBold    = "\033[36;1m"
)

type AppLogHandler struct {
}

func (a AppLogHandler) Trace(v ...interface{}) {
	log.Println(append([]any{Cyan, "[TRACE]", Reset}, v...))
}

func (a AppLogHandler) Debug(v ...interface{}) {
	log.Println(append([]any{BlueBold, "[DEBUG]", Reset}, v...))
}

func (a AppLogHandler) Info(v ...interface{}) {
	log.Println(append([]any{GreenBold, "[INFO]", Reset}, v...))

}

func (a AppLogHandler) Warn(v ...interface{}) {
	log.Println(append([]any{YellowBold, "[WARN]", Reset}, v...))

}

func (a AppLogHandler) Error(v ...interface{}) {
	log.Println(append([]any{RedBold, "[ERROR]", Reset}, v...))

}

func (a AppLogHandler) Fatal(v ...interface{}) {
	log.Println(append([]any{MagentaBold, "[FATAL]", Reset}, v...))

}

func (a AppLogHandler) Panic(v ...interface{}) {
	log.Println(append([]any{MagentaBold, "[PANIC]", Reset}, v...))

}

var AppLogger = &AppLogHandler{}
