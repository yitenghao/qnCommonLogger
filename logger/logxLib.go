package logger

import (
	"fmt"
)

const LOGX_INFO_PREFIX = "UNIQUE-DialogSessionID||"

func LogxPrintf(level int, uid string, format string, parameters ...interface{}) {
	format = fmt.Sprintf("%s\t%s", LOGX_INFO_PREFIX+uid, format)
	switch level {
	case 0:
		LogxDebug(format, parameters...)
	case 1:
		LogxWarnf(format, parameters...)
	case 2:
		LogxErrorf(format, parameters...)
	default:
		LogxDebug(format, parameters...)
	}
}
func FileLogPrintf(name string, level int, uid string, format string, parameters ...interface{}) {
	log := GetFileLoggerMgr().fileWriters[name]
	format = fmt.Sprintf("%s\t%s", LOGX_INFO_PREFIX+uid, format)
	switch level {
	case 0:
		log.Debugf(format, parameters...)
	case 1:
		log.Warnf(format, parameters...)
	case 2:
		log.Errorf(format, parameters...)
	default:
		log.Debugf(format, parameters...)
	}
}
func LogxDebug(format string, parameters ...interface{}) {
	GetCommonFileLogger().Debugf(format, parameters...)
}

func LogxWarnf(format string, parameters ...interface{}) {
	GetCommonFileLogger().Warnf(format, parameters...)
}

func LogxErrorf(format string, parameters ...interface{}) {
	GetCommonFileLogger().Errorf(format, parameters...)
}
