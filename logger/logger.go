package logger

import (
	"log"
)

type logLevel string

const (
	INFO    logLevel = "info"
	ERROR   logLevel = "error"
	WARNING logLevel = "warning"
)

func Info(funcName string, err error, data interface{}) {
	Logger(funcName, INFO, err, data)
}

func Error(funcName string, err error, data interface{}) {
	Logger(funcName, ERROR, err, data)
}

func Warning(funcName string, err error, data interface{}) {
	Logger(funcName, WARNING, err, data)
}

func Logger(funcName string, level logLevel, err error, data interface{}) {
	if err != nil {
		log.Printf("| %s | %s | %s | %v\n", level, funcName, err.Error(), data)
	} else {
		log.Printf("| %s | %s | %s | %v\n", level, funcName, "", data)
	}
}
