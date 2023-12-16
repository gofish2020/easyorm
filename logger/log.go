package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/gofish2020/easyorm/utils"
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
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

type LogLevel int

const (
	LogLevelSlient LogLevel = iota + 1

	LogLevelError
	LogLevelWarn
	LogLevelInfo
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags)

	infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
	warnStr = Yellow + "%s\n" + Reset + Yellow + "[warn] " + Reset
	errStr  = Red + "%s\n" + Reset + Red + "[error] " + Reset

	logLevel = LogLevelInfo
)

func SetLogLevel(l LogLevel) {
	logLevel = l
}
func Infof(msg string, data ...interface{}) {

	if logLevel >= LogLevelInfo {
		logger.Printf(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func Warnf(msg string, data ...interface{}) {
	if logLevel >= LogLevelWarn {
		logger.Printf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}

}
func Errorf(msg string, data ...interface{}) {
	if logLevel >= LogLevelError {
		logger.Printf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func Info(v ...any) {
	if logLevel >= LogLevelInfo {
		logger.Println(append([]interface{}{fmt.Sprintf(infoStr, []interface{}{utils.FileWithLineNum()})}, v...)...)
	}
}

func Warn(v ...any) {
	if logLevel >= LogLevelWarn {
		logger.Println(append([]interface{}{fmt.Sprintf(warnStr, []interface{}{utils.FileWithLineNum()})}, v...)...)
	}
}

func Error(v ...any) {
	if logLevel >= LogLevelError {
		logger.Println(append([]interface{}{fmt.Sprintf(errStr, []interface{}{utils.FileWithLineNum()})}, v...)...)
	}
}
