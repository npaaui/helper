package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Instance *logrus.Logger

// 单个请求统一req_id做链路追踪
func WithReqId(reqId int64) *logrus.Entry {
	return Instance.WithField("req_id", reqId)
}

// 携带code, panic级别的日志可以在middleware.recover中被抓取识别处理接口返回
func WithReqIdAndCode(reqId int64, code int) *logrus.Entry {
	return Instance.WithFields(map[string]interface{}{
		"req_id": reqId,
		"code":   code,
	})
}

func InitLogger(level string, logFolder string) {
	logClient := logrus.New()

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(fmt.Errorf("open log file err: %w", err))
	}
	logClient.Out = src

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic(fmt.Errorf("log level err: %w", err))
	}
	logClient.SetLevel(logLevel)

	// 支持的错误级别
	logLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}

	writeMap := WriterMap{}
	for _, level := range logLevels {
		logWriter, err := rotatelogs.New(
			logFolder+level+"-%Y-%m-%d.log",
			rotatelogs.WithMaxAge(180*24*time.Hour),   // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			panic(fmt.Errorf("new log writer err: %w", err))
		}
		parseLevel, _ := logrus.ParseLevel(level)
		writeMap[parseLevel] = logWriter
	}

	lfHook := NewHook(writeMap, &logrus.JSONFormatter{})
	logClient.AddHook(lfHook)
	Instance = logClient
}
