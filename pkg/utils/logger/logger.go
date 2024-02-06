package logger

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

var (
	standardLogger Logger = &logger{fl: logrus.StandardLogger()}
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

func StandardLogger() Logger {
	return standardLogger
}

func New() Logger {
	return &logger{fl: logrus.New()}
}

func FromLogrus(fl logrus.Ext1FieldLogger) Logger {
	return &logger{fl: fl}
}

type Level = logrus.Level

type Logger interface {
	WithContext(ctx context.Context) Logger

	WithReportCaller(b bool, depth ...uint) Logger

	WithField(key string, value interface{}) Logger
	WithFields(kv ...interface{}) Logger
	WithError(err error) Logger

	SetLevel(level Level) Logger
	WriterLevel(level Level) *io.PipeWriter

	SetOutput(w io.Writer) Logger

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Logr() logr.Logger
	FieldLogger() logrus.FieldLogger
	Logger() *logrus.Logger

	Clone() Logger
}

type logger struct {
	fl           logrus.Ext1FieldLogger
	reportCaller *int
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.withCaller().Tracef(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.withCaller().Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.withCaller().Infof(format, args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.withCaller().Printf(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.withCaller().Warnf(format, args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.withCaller().Warningf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.withCaller().Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.withCaller().Fatalf(format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.withCaller().Panicf(format, args...)
}

func (l *logger) Trace(args ...interface{}) {
	l.withCaller().Trace(args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.withCaller().Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.withCaller().Info(args...)
}

func (l *logger) Print(args ...interface{}) {
	l.withCaller().Print(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.withCaller().Warn(args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.withCaller().Warning(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.withCaller().Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.withCaller().Fatal(args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.withCaller().Panic(args...)
}

func (l *logger) Traceln(args ...interface{}) {
	l.withCaller().Traceln(args...)
}

func (l *logger) Debugln(args ...interface{}) {
	l.withCaller().Debugln(args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.withCaller().Infoln(args...)
}

func (l *logger) Println(args ...interface{}) {
	l.withCaller().Println(args...)
}

func (l *logger) Warnln(args ...interface{}) {
	l.withCaller().Warnln(args...)
}

func (l *logger) Warningln(args ...interface{}) {
	l.withCaller().Warningln(args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.withCaller().Errorln(args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	l.withCaller().Fatalln(args...)
}

func (l *logger) Panicln(args ...interface{}) {
	l.withCaller().Panicln(args...)
}

func (l *logger) WriterLevel(level Level) *io.PipeWriter {
	return l.Logger().WriterLevel(level)
}

func (l *logger) SetLevel(level Level) Logger {
	l.Logger().SetLevel(level)
	return l
}

func (l *logger) WithContext(ctx context.Context) Logger {
	switch t := l.fl.(type) {
	case *logrus.Logger:
		return &logger{fl: t.WithContext(ctx), reportCaller: l.reportCaller}
	case *logrus.Entry:
		return &logger{fl: t.WithContext(ctx), reportCaller: l.reportCaller}
	}
	panic(fmt.Sprintf("unexpected logger type %T", l.fl))
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{fl: l.fl.WithField(key, value), reportCaller: l.reportCaller}
}

func (l *logger) WithFields(kv ...interface{}) Logger {
	log := &logger{fl: l.fl}
	for i := 0; i < len(kv); i += 2 {
		log = &logger{fl: log.fl.WithField(fmt.Sprintf("%v", kv[i]), kv[i+1]), reportCaller: l.reportCaller}
	}
	return log
}

func (l *logger) WithError(err error) Logger {
	return &logger{fl: l.fl.WithError(err), reportCaller: l.reportCaller}
}

func (l *logger) WithReportCaller(b bool, depth ...uint) Logger {
	if !b {
		return &logger{fl: l.fl}
	}
	var d int
	if len(depth) > 0 {
		d = int(depth[0])
	} else {
		d = 0
	}
	return &logger{fl: l.fl, reportCaller: &d}
}

func (l *logger) Logr() logr.Logger {
	return logrusr.New(l.fl)
}

func (l *logger) FieldLogger() logrus.FieldLogger {
	return l.fl
}

func (l *logger) Logger() *logrus.Logger {
	switch t := l.fl.(type) {
	case *logrus.Logger:
		return t
	case *logrus.Entry:
		return t.Logger
	}
	panic(fmt.Sprintf("unexpected logger type %T", l.fl))
}

func (l *logger) SetOutput(w io.Writer) Logger {
	l.Logger().SetOutput(w)
	return l
}

func (l *logger) Clone() Logger {
	n := logrus.New()
	switch t := l.fl.(type) {
	case *logrus.Logger:
		n.Level = t.Level
		n.Out = t.Out
		n.Formatter = t.Formatter
		n.Hooks = t.Hooks
		return &logger{fl: n, reportCaller: l.reportCaller}
	case *logrus.Entry:
		t = t.Dup()
		n.Level = t.Logger.Level
		n.Out = t.Logger.Out
		n.Formatter = t.Logger.Formatter
		n.Hooks = t.Logger.Hooks
		t.Logger = n
		return &logger{fl: t, reportCaller: l.reportCaller}
	}
	panic(fmt.Sprintf("unexpected logger type %T", l.fl))
}

func (l *logger) withCaller() logrus.Ext1FieldLogger {
	if l.reportCaller == nil {
		return l.fl
	}
	pcs := make([]uintptr, 1)
	runtime.Callers(3+*l.reportCaller, pcs)
	f, _ := runtime.CallersFrames(pcs).Next()
	pkg := getPackageName(f.Function)
	return l.fl.WithField("caller", fmt.Sprintf("%s/%s:%d", pkg, filepath.Base(f.File), f.Line)).WithField("func", f.Func.Name())
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be a better way...
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
