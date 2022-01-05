package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"qnCommonLogger/common"
)

func setupLoggerCommonByDateRotate(log *logrus.Logger, level logrus.Level, fileName, logDir string, maxSizeM, rotateMaxBackups, rotateMaxAge, skip int, report bool) {
	if log == nil {
		common.CriticalExitf("get empty logger pointer\n")
		return
	}
	isExist, err := common.PathExists(logDir)
	if err != nil {
		common.CriticalExitf("failed to create trace output file: %v", err)
	}
	if !isExist {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			common.CriticalExitf("failed to create trace output dir: %v", err)
		}
	}

	path := logDir + "/" + fileName
	writer := &lumberjack.Logger{
		Filename:   path,             // 日志文件路径
		MaxSize:    maxSizeM,         // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: rotateMaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     rotateMaxAge,     // 文件最多保存多少天
		Compress:   false,            // 是否压缩
	}
	// 设置日志格式为json格式
	log.SetFormatter(&logrus.JSONFormatter{})
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(writer)
	log.SetReportCaller(false)
	// log.SetOutput(os.Stdout)
	// 设置日志级别为info以上
	if report {
		log.AddHook(&DateLogHook{skip: skip})
	}

	log.SetLevel(level)

	GetFileLoggerMgr().AddOne(commonLog, Logger{Logger: log})
}
