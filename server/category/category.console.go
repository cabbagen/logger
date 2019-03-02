package category

import (
	"fmt"
	"time"
	"strings"
)

type ConsoleLogger struct {
	LoggerStruct
}

var (
	white      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	green      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	yello      = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	blue       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	red        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	reset      = string([]byte{27, 91, 48, 109})
)

var colorLevelMap map[string]string = map[string]string {
	"ALL": white,
	"INFO": green,
	"WARN": yello,
	"DEBUG": blue,
	"ERROR": red,
}

func NewConsoleLogger(appid, level string) ConsoleLogger {
	consoleLogger := ConsoleLogger{}
	
	consoleLogger.Appid = appid
	consoleLogger.Level = level
	consoleLogger.IsCache = false
	consoleLogger.Cacher = nil
	
	return consoleLogger
}

// 根据不同的级别，打印出不同的颜色
func (cl ConsoleLogger) Log(message string) {
	
	prefix := cl.GetTimePrefix()
	
	color := colorLevelMap[strings.ToUpper(cl.Level)]
	
	if color == "" {
		fmt.Printf("%s app[%s] [%s]: %s\n", prefix, cl.Appid, cl.Level, message)
		return
	}
	
	fmt.Printf("%s %s app[%s] [%s]: %s %s\n", color, prefix, cl.Appid, cl.Level, message, reset)
}

func (cl ConsoleLogger) GetTimePrefix() string {
	now := time.Now()
	
	return fmt.Sprintf("%d-%d-%d %d/%d/%d ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

