package candy

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/valyala/bytebufferpool"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*

// 设置日志 目录、级别、日志最大保留天数。 (不设置则所有日志打印到控制台)
candy.SetLog("./log", "debug", 30)

// 打印日志到 default.log
candy.Debug("hello", "world")

// 自定义日志文件名
l := candy.NewLog("test.log")
l.Debug("hello", "world")
// 打印错误时，自动追踪链路，并额外记录在 error.log
l.Error(error)

*/

type level uint8

const (
	levelDebug level = iota + 1
	levelInfo
	levelWarn
	levelError
)

var (
	logDir    string = ""         // 日志目录
	logLevel  level  = levelDebug // 日志级别
	logMaxAge uint8  = 30         // 日志保存天数

	logBufferPool = &bytebufferpool.Pool{} // buffer池

	levelAlias = [...]string{"", "[DEBUG]", "[INFO]", "[WARN]", "[ERROR]"}

	defaultLogger *Logger = NewLogger("")
	errorLogger   *Logger = NewLogger("")
)

// 设置日志目录、级别、最大保留天数
func SetLog(dir string, level string, maxAge uint8) {
	logDir = strings.TrimSuffix(dir, "/") + "/"
	logMaxAge = maxAge
	switch strings.ToLower(level) {
	case "debug":
		logLevel = levelDebug
	case "info":
		logLevel = levelInfo
	case "warn":
		logLevel = levelWarn
	case "error":
		logLevel = levelError
	}

	defaultLogger = NewLogger("default.log")
	errorLogger = NewLogger("error.log")
}

func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

func Info(args ...any) {
	defaultLogger.Info(args...)
}

func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

func Error(args ...any) {
	defaultLogger.Error(args...)
}

type Logger struct {
	io.Writer
}

func NewLogger(filename string) *Logger {
	var w io.Writer = os.Stdout
	if logDir != "" {
		w = NewWriter(logDir + filename)
	}
	return &Logger{w}
}

func NewWriter(filename string) io.Writer {
	if filename == "" {
		return os.Stdout
	}

	return &lumberjack.Logger{
		Filename:  filename,
		MaxSize:   0,              // 最大日志文件大小，超过该大小会进行日志文件切割
		MaxAge:    int(logMaxAge), // 最大保存天数，超过该天数的日志文件会被删除
		LocalTime: true,
		Compress:  false,
	}
}

func (l *Logger) Debug(args ...any) {
	write(l, levelDebug, args...)
}

func (l *Logger) Info(args ...any) {
	write(l, levelInfo, args...)
}

func (l *Logger) Warn(args ...any) {
	write(l, levelWarn, args...)
}

func (l *Logger) Error(args ...any) {
	write(l, levelError, args...)
}

func write(l *Logger, level level, args ...any) {
	if level < logLevel {
		return
	}

	if len(args) == 0 {
		return
	}

	buf := logBufferPool.Get()
	defer logBufferPool.Put(buf)
	buf.WriteString(time.Now().Format(time.DateTime))
	buf.WriteString(" " + levelAlias[level] + " ")

	for _, arg := range args {
		switch v := arg.(type) {
		case []byte:
			buf.Write(v)
		case string:
			buf.WriteString(v)
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, bool:
			buf.WriteString(fmt.Sprintf("%v", v))
		case error:
			buf.WriteString(v.Error())
		default:
			b, err := json.Marshal(arg)
			if err != nil {
				buf.WriteString(fmt.Sprintf("%v", v))
			} else {
				buf.Write(b)
			}
		}
		buf.WriteByte(' ')
	}

	buf.WriteByte('\n')
	l.Write(buf.Bytes())

	if level == levelError {
		buf.Bytes()[len(buf.Bytes())-1] = ' '
		b, _ := json.Marshal(Callers(2)[3:])
		buf.Write(b)
		buf.WriteByte('\n')
		errorLogger.Write(buf.Bytes())
	}
}
