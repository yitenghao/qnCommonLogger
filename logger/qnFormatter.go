package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"qnCommonLogger/common"
)

type QnFormatter struct {
}
type standardLogFormat struct {
	Service string      `json:"service"`
	From    interface{} `json:"from,omitempty"`
	Out     interface{} `json:"out,omitempty"`
	Time    float32     `json:"time,omitempty"`
}

func (s *QnFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05.999")
	if !strings.HasSuffix(entry.Message, "\n") {
		entry.Message = entry.Message + "\n"
	}
	if entry.Caller == nil {
		msg := fmt.Sprintf("%s [%s]\t%s", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
		return []byte(msg), nil
	}

	file, line, fName := entry.Caller.File, entry.Caller.Line, entry.Caller.Function

	msg := fmt.Sprintf("%s [%s] [%s:%d %s]\t%s", timestamp, strings.ToUpper(entry.Level.String()), file, line, fName, entry.Message)
	return []byte(msg), nil
}

type StandardLog struct {
	ServiceName string      `json:"serviceName"`
	Msg         interface{} `json:"msg"`
	Time        int64       `json:"time"`
}

func ToStandardFormatEx(isInput bool, data StandardLog) string {
	s := standardLogFormat{
		Service: data.ServiceName,
	}
	if isInput {
		s.From = data.Msg
	} else {
		s.Out = data.Msg
		s.Time = float32(data.Time) / float32(time.Second)
	}
	bs, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("marshal data %+v, err: %s\n", s, err.Error())
	}
	str := string(bs)
	return str
}

func setupDIYFormatByDate(logs *logrus.Logger, level logrus.Level, fileName, logDir, serviceName string, rotateMaxAge, skip int, report bool) {
	if logs == nil {
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

	dateStr := time.Now().Format("20060102")
	logFullPath := logDir + "/" + fileName + "." + dateStr
	// os.O_O_CREATE auto create
	logFile, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		common.CriticalExitf("Open logger file fail: %s", err.Error())
		return
	}
	fileWriter := LogFileWriterByDate{
		dirPath:      logDir,
		fileName:     fileName,
		file:         logFile,
		dateStr:      dateStr,
		mutex:        &sync.Mutex{},
		rotateMaxAge: rotateMaxAge,
	}
	fileWriter.mustDoMillFirst()
	// 设置日志格式为json格式
	logs.SetFormatter(&QnFormatter{})
	// 开启调用上下文记录
	logs.SetReportCaller(false)
	// 日志消息输出可以是任意的io.writer类型
	logs.SetOutput(&fileWriter)
	if report {
		logs.AddHook(&DateLogHook{skip: skip})
	}

	logs.SetLevel(level)

	GetFileLoggerMgr().AddOne(qnLog, Logger{
		Logger:      logs,
		ServiceName: serviceName,
	})
}
