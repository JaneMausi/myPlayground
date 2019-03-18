package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	//LOG_EMERG = iota
	//LOG_ALERT
	//LOG_CRIT
	LOG_ERR = iota
	//LOG_WARNING
	//LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
	LOG_NUM
)

var infoLogger *log.Logger
var errLogger *log.Logger
var debugLogger *log.Logger

var loggerMap = [LOG_NUM]struct {
	level  string
	path   string
	logger *log.Logger
}{
	{"LOG_ERR", "", errLogger},
	{"LOG_INFO", "", infoLogger},
	{"LOG_DEBUG", "", debugLogger},
}

func InitLogWriter() {
	logfiles := []string {
		"error.log",
		"info.log",
		"debug.log",
	}
	for index := 0; index < len(loggerMap); index++ {
		loggerMap[index].path = logfiles[index]
		if index == LOG_DEBUG && loggerMap[index].path == "" {
			loggerMap[index].logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
			continue
		}
		if loggerMap[index].path == "" {
			Fatal(fmt.Sprintf("Log file for %s not defined, check your regcoin.conf file please.",
				loggerMap[index].level))
		} else {
			logWriter, err := os.OpenFile(loggerMap[index].path,
				os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				Fatal(err)
			}
			loggerMap[index].logger = log.New(logWriter, "", log.LstdFlags)
		}
	}
}

func WriteLog(level int, msg ...interface{}) {
	_, filepath, line, _ := runtime.Caller(1)
	filename := path.Base(filepath)

	var paramSlice []string
	paramSlice = append(paramSlice, filename, fmt.Sprint(line))
	for _, param := range msg {
		paramSlice = append(paramSlice, fmt.Sprint(param))
	}
	loggerMap[level].logger.Println(strings.Join(paramSlice, ", "))
}

func Fatal(msg ...interface{}) {
	_, filepath, line, _ := runtime.Caller(1)
	filename := path.Base(filepath)

	var paramSlice []string
	paramSlice = append(paramSlice, filename, fmt.Sprint(line))
	for _, param := range msg {
		paramSlice = append(paramSlice, fmt.Sprint(param))
	}
	log.Fatal(strings.Join(paramSlice, ", "))
}
