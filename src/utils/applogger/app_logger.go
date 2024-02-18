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

var logFlags = log.LstdFlags | log.LUTC | log.Lmsgprefix | log.Lshortfile
var (
	traceLogger = log.New(os.Stdout, _cyanbold+"[TRACE] "+_reset, logFlags)
	debugLogger = log.New(os.Stdout, _bluebold+"[DEBUG] "+_reset, logFlags)
	infoLogger  = log.New(os.Stdout, _greenbold+"[INFO] "+_reset, logFlags)
	warnLogger  = log.New(os.Stdout, _yellowbold+"[WARN] "+_reset, logFlags)
	errorLogger = log.New(os.Stderr, _redbold+"[ERROR] "+_reset, logFlags)
	fatalLogger = log.New(os.Stderr, _magentabold+"[FATAL] "+_reset, logFlags)
	panicLogger = log.New(os.Stderr, _magentabold+"[PANIC] "+_reset, logFlags)
)

func Trace(v ...interface{}) {
	traceLogger.Println(v...)
}

func Debug(v ...interface{}) {
	debugLogger.Println(v...)
}

func Info(v ...interface{}) {
	infoLogger.Println(v...)

}

func Warn(v ...interface{}) {
	warnLogger.Println(v...)

}

func Error(v ...interface{}) {
	errorLogger.Println(v...)

}

func Fatal(v ...interface{}) {
	fatalLogger.Println(v...)
	os.Exit(1)
}

func Panic(v ...interface{}) {
	panicLogger.Println(v...)
	panic(v)
}
