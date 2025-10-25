package logging

import (
	"MediaWarp/internal/config"
	"io"

	"github.com/sirupsen/logrus"
)

var (
	accessLogger  = logrus.New() // 访问日志
	serviceLogger = logrus.New() // 服务日志
)

func init() {
	accessLogger.SetFormatter(&LoggerAccessFormatter{})
	serviceLogger.SetFormatter(&LoggerServiceFormatter{})
}

func Init() {
	serviceLogger.SetReportCaller(false) // 关闭报告调用方

	if !config.Logger.AccessLogger.Console { // 访问日志不输出到终端
		accessLogger.Out = io.Discard
	}

	if !config.Logger.ServiceLogger.Console { // 服务日志不输出到终端
		serviceLogger.Out = io.Discard
	}

	if config.Logger.AccessLogger.File {
		accessLogger.AddHook(NewLoggerFileHook(false))
	}

	if config.Logger.ServiceLogger.File {
		serviceLogger.AddHook(NewLoggerFileHook(true))
	}
}

// 访问日志
//
// 默认日志级别为 Info
func AccessLog(format string, args ...any) {
	accessLogger.Infof(format, args...)
}

func AccessDebugf(format string, args ...any) {
	accessLogger.Debugf(format, args...)
}

func AccessWarningf(format string, args ...any) {
	accessLogger.Warningf(format, args...)
}

// 服务日志
//
// Debug 级别日志
func Debug(args ...any) {
	serviceLogger.Debug(args...)
}

func Debugf(format string, args ...any) {
	serviceLogger.Debugf(format, args...)
}

// 服务日志
//
// Info 级别日志
func Info(args ...any) {
	serviceLogger.Info(args...)
}

func Infof(format string, args ...any) {
	serviceLogger.Infof(format, args...)
}

// 服务日志
//
// Warning 级别日志
func Warning(args ...any) {
	serviceLogger.Warning(args...)
}

func Warningf(format string, args ...any) {
	serviceLogger.Warningf(format, args...)
}

// 服务日志
//
// Error 级别日志
func Error(args ...any) {
	serviceLogger.Error(args...)
}

func Errorf(format string, args ...any) {
	serviceLogger.Errorf(format, args...)
}

// 服务日志
//
// 设置日志级别
func SetLevel(level logrus.Level) {
	serviceLogger.SetLevel(level)
}
