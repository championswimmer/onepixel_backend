package applogger

import (
	"log"
	"os"
)

const (
	_reset       = "\033[0m"
	_red         = "\033[31m"
	_green       = "\033[32m"
	_yellow      = "\033[33m"
	_blue        = "\033[34m"
	_magenta     = "\033[35m"
	_cyan        = "\033[36m"
	_white       = "\033[37m"
	_redbold     = "\033[31;1m"
	_greenbold   = "\033[32;1m"
	_yellowbold  = "\033[33;1m"
	_bluebold    = "\033[34;1m"
	_magentabold = "\033[35;1m"
	_cyanbold    = "\033[36;1m"
)

func Trace(v ...interface{}) {
	log.Println(append([]any{_cyan, "[TRACE]", _reset}, v...))
}

func Debug(v ...interface{}) {
	log.Println(append([]any{_bluebold, "[DEBUG]", _reset}, v...))
}

func Info(v ...interface{}) {
	log.Println(append([]any{_greenbold, "[INFO]", _reset}, v...))

}

func Warn(v ...interface{}) {
	log.Println(append([]any{_yellowbold, "[WARN]", _reset}, v...))

}

func Error(v ...interface{}) {
	log.Println(append([]any{_redbold, "[ERROR]", _reset}, v...))

}

func Fatal(v ...interface{}) {
	log.Println(append([]any{_magentabold, "[FATAL]", _reset}, v...))
	os.Exit(1)
}

func Panic(v ...interface{}) {
	log.Println(append([]any{_magentabold, "[PANIC]", _reset}, v...))
	panic(v)
}
