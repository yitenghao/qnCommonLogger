package logger

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/linnv/logx"
	"github.com/sirupsen/logrus"
)

const (
	logSuccessSessionFileName = "session_success.log"
	logFailSessionFileName    = "session_fail.log"

	commonLog = "common"
	qnLog     = "qnApi"
)

var (
	fileLoggerMgr *FileLoggerMgr
	once          sync.Once
)

func GetFileLoggerMgr() *FileLoggerMgr {
	once.Do(func() {
		fileLoggerMgr = NewLogFileMgr()
	})

	return fileLoggerMgr
}

func NewLogFileMgr() *FileLoggerMgr {
	fm := make(map[string]Logger)
	return &FileLoggerMgr{
		fileWriters: fm,
	}
}

func (f *FileLoggerMgr) AddOne(key string, value Logger) {
	if f == nil {
		logx.Errorf("LogFileMgr empty\n")
		return
	}

	if f.fileWriters == nil {
		f.fileWriters = make(map[string]Logger)
	}

	f.fileWriters[key] = value
	return
}

func GetCommonFileLogger() Logger {
	return GetFileLoggerMgr().getCommonFileLogger()
}

func GetSuccessFileLogger() Logger {
	return GetFileLoggerMgr().getSuccessFileLogger()
}

func GetFailFileLogger() Logger {
	return GetFileLoggerMgr().getFailFileLogger()
}
func GetQnFileLogger() Logger {
	return GetFileLoggerMgr().getQnFileLogger()
}
func (f *FileLoggerMgr) getCommonFileLogger() Logger {
	return f.fileWriters[commonLog]
}

func (f *FileLoggerMgr) getSuccessFileLogger() Logger {
	return f.fileWriters[logSuccessSessionFileName]
}

func (f *FileLoggerMgr) getFailFileLogger() Logger {
	return f.fileWriters[logFailSessionFileName]
}
func (f *FileLoggerMgr) getQnFileLogger() Logger {
	return f.fileWriters[qnLog]
}

type Logger struct {
	*logrus.Logger
	ServiceName string
}

func (l Logger) ParseQnApiLogFormat(isInput bool, t int64, data interface{}) {
	s := standardLogFormat{
		Service: l.ServiceName,
	}
	if isInput {
		s.From = data
	} else {
		s.Out = data
		s.Time = float32(t) / float32(time.Second)
	}
	bs, err := json.Marshal(s)
	if err != nil {
		LogxPrintf(2, "", "ParseQnApiLogFormat Marshal err:%s", err.Error())
	}
	str := string(bs)
	l.Info(str)
}
