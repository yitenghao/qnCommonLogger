package logger

type FileLoggerMgr struct {
	fileWriters map[string]Logger // key: fileName value: LogFileWriter
}
