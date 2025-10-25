package logging

import (
	"MediaWarp/constants"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerServiceFormatter struct{}

func (l *LoggerServiceFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据日志级别设置颜色
	var colorCode uint8
	switch entry.Level {
	case logrus.DebugLevel:
		colorCode = constants.ColorBlue
	case logrus.InfoLevel:
		colorCode = constants.ColorGreen
	case logrus.WarnLevel:
		colorCode = constants.ColorYellow
	case logrus.ErrorLevel:
		colorCode = constants.ColorRed
	default:
		colorCode = constants.ColorGray
	}

	// 设置文本Buffer
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}
	// 时间格式化
	formatTime := entry.Time.Format(time.DateTime)

	fmt.Fprintf(
		b,
		"\033[3%dm【%s】\033[0m\t%s | %s\n", // 长度需要算是上控制字符的长度
		colorCode,
		strings.ToUpper(entry.Level.String()),
		formatTime,
		entry.Message,
	)
	return b.Bytes(), nil
}

var _ logrus.Formatter = (*LoggerServiceFormatter)(nil)

type LoggerAccessFormatter struct{}

// 实现Format方法
func (l *LoggerAccessFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}

	fmt.Fprint(
		b,
		entry.Message+"\n",
	)
	return b.Bytes(), nil
}

var _ logrus.Formatter = (*LoggerAccessFormatter)(nil)
