package logger

import (
	"os"
	"path"
	"strconv"

	"github.com/linnv/logx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"qnCommonLogger/common"
)

type LogFileWriter struct {
	dirPath  string
	fileName string
	maxSize  int64
	file     *os.File
	size     int64
}

func (p *LogFileWriter) Write(data []byte) (int, error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	n, err := p.file.Write(data)
	p.size += int64(n)
	// 文件最大 64 K byte
	if p.size > p.maxSize {
		p.file.Close()
		logx.Debugln("log file full")
		count, err := common.CountDirFileNum(p.dirPath)
		if err != nil {
			logx.Warnf("CountDirFileNum err: %s\n", err.Error())
			return n, err
		}
		fullPath := p.dirPath + "/" + p.fileName
		renamePath := fullPath + "." + strconv.Itoa(count)
		if err := os.Rename(fullPath, renamePath); err != nil {
			// 出错则不切割
			logx.Errorf("Rename file: %s, to %s, err: %s\n", fullPath, renamePath, err.Error())
			return n, err
		}
		p.file, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			logx.Errorf("OpenFile a new file err: %s\n", err.Error())
			return n, err
		}
		p.size = 0
	}
	return n, err
}

func setupLogger(log *logrus.Logger, fileName string, logDir string, sizeLimit int64, level logrus.Level) {
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

	logFullPath := logDir + "/" + fileName
	// os.O_O_CREATE auto create
	logFile, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		common.CriticalExitf("Open logger file fail: %s", err.Error())
		return
	}
	info, err := logFile.Stat()
	if err != nil {
		logx.Errorf("get log file state err: %s\n", err.Error())
		common.CriticalExitf(err.Error())
		return
	}
	fileWriter := LogFileWriter{logDir, fileName, sizeLimit, logFile, info.Size()}
	// 设置日志格式为json格式
	log.SetFormatter(&logrus.JSONFormatter{})
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(&fileWriter)
	// log.SetOutput(os.Stdout)
	// 设置日志级别为info以上
	log.SetLevel(level)

	GetFileLoggerMgr().AddOne(commonLog, Logger{Logger: log})
}

// 按文件大小分割
func SetupLoggerCommon(logDir, logName string, sizeLimit int64, level logrus.Level) {
	if logDir == "" {
		configFilePath := viper.GetString("c")
		if configFilePath != "" {
			logDir = path.Join(configFilePath, "../log")
		} else {
			logDir = path.Join(configFilePath, "./log")
		}
	}
	logx.Debugf("SetupLoggerCommon using log dir : [%s]\n", logDir)

	log := logrus.New()
	setupLogger(log, logName, logDir, sizeLimit, level)
}

// 按时间分割
func SetupLoggerCommonByDate(logDir, logName string, rotateMaxAge, skip int, report bool, level logrus.Level) {
	if logDir == "" {
		configFilePath := viper.GetString("c")
		if configFilePath != "" {
			logDir = path.Join(configFilePath, "../log")
		} else {
			logDir = path.Join(configFilePath, "./log")
		}
	}
	logx.Debugf("SetupLoggerCommonByDate using log dir : [%s]\n", logDir)

	log := logrus.New()
	setupLoggerCommonByDate(log, level, logName, logDir, rotateMaxAge, skip, report)
}

/*
	按文件大小和时间并备份分割文件
	logDir:日志所在目录
	maxSizeM:日志最大容量 MB
	rotateMaxBackups:最多保存备份数
	rotateMaxAge:文件最长保存天数
	level:日志等级
*/
func SetupLoggerCommonRotate(logDir, logName string, maxSizeM, rotateMaxBackups, rotateMaxAge, skip int, report bool, level logrus.Level) {
	if logDir == "" {
		configFilePath := viper.GetString("c")
		if configFilePath != "" {
			logDir = path.Join(configFilePath, "../log")
		} else {
			logDir = path.Join(configFilePath, "./log")
		}
	}
	logx.Debugf("SetupLoggerCommon using log dir : [%s]\n", logDir)

	log := logrus.New()
	setupLoggerCommonByDateRotate(log, level, logName, logDir, maxSizeM, rotateMaxBackups, rotateMaxAge, skip, report)
}

// qnzs-go api输入输出标准日志
// 按时间分割
// 输出json
func SetupQnFormatByDate(logDir, logName, serviceName string, rotateMaxAge, skip int, report bool, level logrus.Level) {
	if logDir == "" {
		configFilePath := viper.GetString("c")
		if configFilePath != "" {
			logDir = path.Join(configFilePath, "../logQn")
		} else {
			logDir = path.Join(configFilePath, "./logQn")
		}
	}
	logx.Debugf("SetupQnFormatByDate using log dir : [%s]\n", logDir)

	log := logrus.New()
	setupDIYFormatByDate(log, level, logName, logDir, serviceName, rotateMaxAge, skip, report)
}

func parseLogLevel(level int) logrus.Level {
	switch level {
	case 0:
		return logrus.DebugLevel
	case 1:
		return logrus.WarnLevel
	case 2:
		return logrus.ErrorLevel
	}
	return logrus.DebugLevel
}
