package main

import (
	"github.com/sirupsen/logrus"

	"qnCommonLogger/logger"
)

const (
	DEBUG_LEVEL = 0
	WARN_LEVEL  = 1
	ERROR_LEVEL = 2
	REDIS_ERROR = -1
)

// example
func main() {
	logger.SetupLoggerCommonByDate("log", "demo.log", 7, 0, false, logrus.DebugLevel)
	logger.SetupQnFormatByDate("log", "qnApi.log", "demoService", 7, 0, false, logrus.DebugLevel)

	logger.LogxPrintf(DEBUG_LEVEL, "", "dsfsd")
	logger.GetQnFileLogger().ParseQnApiLogFormat(true, 0, []string{"12345", "aaaaa", "-----"})
	logger.GetQnFileLogger().ParseQnApiLogFormat(false, 100, "success")

}
