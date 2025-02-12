package logger

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gookit/color"
)

type LogLevel int

const (
	LogLevelError         LogLevel = iota + 1 // 错误信息
	LogLevelWarn                              // 警告信息
	LogLevelInfo                              // 一般信息
	LogLevelDebug                             // 调试信息
	LogLevelVerbose                           // 最详细的日志级别
	LogLevelDebugResponse                     // 开发专用级别

)

var (
	Red         = color.Red.Render
	Cyan        = color.Cyan.Render
	Yellow      = color.Yellow.Render
	White       = color.White.Render
	Blue        = color.Blue.Render
	Purple      = color.Style{color.Magenta, color.OpBold}.Render
	LightRed    = color.Style{color.Red, color.OpBold}.Render
	LightGreen  = color.Style{color.Green, color.OpBold}.Render
	LightWhite  = color.Style{color.White, color.OpBold}.Render
	LightCyan   = color.Style{color.Cyan, color.OpBold}.Render
	LightYellow = color.Style{color.Yellow, color.OpBold}.Render
	LightBlue   = color.Style{color.Blue, color.OpBold}.Render
)

func init() {
	// 设置日志输出为0，不显示时间戳
	log.SetFlags(0)
}

type Logger struct {
	Level LogLevel
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{Level: level}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.Level >= LogLevelDebug {
		log.Printf("[%s] [%s] %s", Cyan(GetTime()), LightRed("ERROR"), fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.Level >= LogLevelInfo {
		log.Printf("[%s] [%s] %s", Cyan(GetTime()), LightGreen("INFO"), fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.Level >= LogLevelWarn {
		log.Printf("[%s] [%s] %s", Cyan(GetTime()), LightYellow("WARN"), fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.Level >= LogLevelDebug {
		log.Printf("[%s] [%s] %s", Cyan(GetTime()), LightBlue("DEBUG"), fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Verbosef(format string, v ...interface{}) {
	if l.Level >= LogLevelVerbose {
		log.Printf("[%s] [%s] %s", Cyan(GetTime()), LightCyan("VERBOSE"), fmt.Sprintf(format, v...))
	}
}

func (l *Logger) DebugResponsef(format string, v ...interface{}) {
	if l.Level == LogLevelDebugResponse {
		log.Printf("[%s] [%s] %s", Cyan(GetTime()), LightRed("DEBUG_RESPONSE"), fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Success(tags []string, url string, title string) {
	// 使用strings.Builder来优化字符串拼接性能
	var sb strings.Builder
	for _, tag := range tags {
		sb.WriteString("[")
		sb.WriteString(Yellow(tag))
		sb.WriteString("]")
	}
	// 添加日志级别检查,保持与其他日志方法一致
	log.Printf("[%s] [%s] [%s] %s %s [%s]", Cyan(GetTime()), LightGreen("SUCCESS"), Cyan("TCP/HTTP"), sb.String(), Blue(url), Blue(title))
}

func GetTime() string {
	return time.Now().Format("15:04:05")
}
