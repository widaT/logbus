package logbus

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var (
	logPackage         string
	maximumCallerDepth int = 10
	minimumCallerDepth     = 1
	loggers                = make(map[string]*Logger)
	defaultLogger          = NewLogger(DebugLevel, "default")
)

func NewLogger(level Level, prefix string) *Logger {
	if logger, found := loggers[prefix]; found {
		return logger
	}
	l := &Logger{
		level:  level,
		prefix: prefix,
		output: os.Stdout,
	}
	loggers[prefix] = l
	return l
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func getCaller() *runtime.Frame {
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		if !strings.Contains(pkg, "logbus") && pkg != logPackage {
			return &f
		}
	}
	return nil
}
func Info(v ...interface{}) { defaultLogger.Info(v...) }

func Trace(v ...interface{}) { defaultLogger.Trace(v...) }

func Debug(v ...interface{}) { defaultLogger.Debug(v...) }

func Warn(v ...interface{}) { defaultLogger.Warn(v...) }

func Error(v ...interface{}) { defaultLogger.Error(v...) }

func Panic(v ...interface{}) { defaultLogger.Panic(v...) }

func Infof(format string, v ...interface{}) { defaultLogger.Infof(format, v...) }

func Tracef(format string, v ...interface{}) { defaultLogger.Tracef(format, v...) }

func Debugf(format string, v ...interface{}) { defaultLogger.Debugf(format, v...) }

func Warnf(format string, v ...interface{}) { defaultLogger.Warnf(format, v...) }

func Errorf(format string, v ...interface{}) { defaultLogger.Errorf(format, v...) }

func Panicf(format string, v ...interface{}) { defaultLogger.Panicf(format, v...) }

func SetLogLevel(prefix string, level Level) error {
	if l, found := loggers[prefix]; found {
		l.level = level
		l.SetLevel(level)
		return nil
	}
	return fmt.Errorf("logger [%v] not found", prefix)
}

func SetDefaultLogLevel(level Level) error {
	defaultLogger.SetLevel(level)
	return nil
}

func GetLoggers() map[string]*Logger {
	return loggers
}

func StringToLevel(level string) Level {
	l := DebugLevel
	switch level {
	case "trace":
		l = TraceLevel
	case "debug":
		l = DebugLevel
	case "info":
		l = InfoLevel
	case "warn":
		l = WarnLevel
	case "error":
		l = ErrorLevel
	}
	return l
}

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	}
	return "UNKOWN"
}

type Logger struct {
	level  Level
	prefix string
	output io.Writer
}

func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Info(v ...interface{}) {
	if l.level < InfoLevel {
		return
	}
	l.write(InfoLevel, nil, v...)
}

func (l *Logger) Trace(v ...interface{}) {
	if l.level < TraceLevel {
		return
	}
	l.write(TraceLevel, nil, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level < DebugLevel {
		return
	}
	l.write(DebugLevel, nil, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level < WarnLevel {
		return
	}
	l.write(WarnLevel, nil, v...)
}

func (l *Logger) Error(v ...interface{}) {
	if l.level < ErrorLevel {
		return
	}
	l.write(ErrorLevel, nil, v...)
}

func (l *Logger) Panic(v ...interface{}) {
	l.write(PanicLevel, nil, v...)
	os.Exit(-1)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level < InfoLevel {
		return
	}
	l.write(InfoLevel, &format, v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if l.level < TraceLevel {
		return
	}
	l.write(TraceLevel, &format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level < DebugLevel {
		return
	}
	l.write(DebugLevel, &format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level < WarnLevel {
		return
	}
	l.write(WarnLevel, &format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level < ErrorLevel {
		return
	}
	l.write(ErrorLevel, &format, v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.write(PanicLevel, &format, v...)
	os.Exit(-1)
}

func (l *Logger) write(level Level, format *string, v ...interface{}) {
	frame := getCaller()
	temp := fmt.Sprintf("[%s] %5s %s: [%s:%d] [%s]",
		time.Now().Format("2006-01-02 15:04:05.000"),
		level,
		l.prefix,
		path.Base(frame.File),
		frame.Line,
		frame.Function,
	)
	if format == nil {
		fmt.Fprintf(l.output, "%s %s\n", temp, fmt.Sprint(v...))
		return
	}
	fmt.Fprintf(l.output, "%s %s\n", temp, fmt.Sprintf(*format, v...))
}
