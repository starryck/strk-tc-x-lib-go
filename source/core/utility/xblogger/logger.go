package xblogger

import (
	"fmt"
	"maps"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/starryck/strk-tc-x-lib-go/source/core/base/xbcfg"
	"github.com/starryck/strk-tc-x-lib-go/source/core/toolkit/xbslice"
	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xbctnr"
	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xberror"
	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xbjson"
)

var mLogger *Logger

var (
	SkipKey  = "#skip"
	PanicKey = "panic"
	ErrorKey = logrus.ErrorKey
)

type (
	Logger = logrus.Logger
	Entry  = logrus.Entry
	Level  = logrus.Level
	Fields = logrus.Fields
)

func GetLogger() *Logger {
	if mLogger == nil {
		mLogger = newLogger()
	}
	return mLogger
}

func newLogger() *Logger {
	logger := (&loggerBuilder{}).
		initialize().
		setSeverity().
		setFormatter().
		build()
	return logger
}

func Trace(args ...any) {
	GetLogger().Trace(args...)
}

func Tracef(format string, args ...any) {
	GetLogger().Tracef(format, args...)
}

func Debug(args ...any) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...any) {
	GetLogger().Debugf(format, args...)
}

func Info(args ...any) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...any) {
	GetLogger().Infof(format, args...)
}

func Warn(args ...any) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...any) {
	GetLogger().Warnf(format, args...)
}

func Error(args ...any) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...any) {
	GetLogger().Errorf(format, args...)
}

func Fatal(args ...any) {
	GetLogger().Fatal(args...)
}

func Fatalf(format string, args ...any) {
	GetLogger().Fatalf(format, args...)
}

func Panic(args ...any) {
	GetLogger().Panic(args...)
}

func Panicf(format string, args ...any) {
	GetLogger().Panicf(format, args...)
}

func WithError(err error) *Entry {
	return GetLogger().WithError(err)
}

func WithFields(fields Fields) *Entry {
	return GetLogger().WithFields(fields)
}

func GetLevel() Level {
	return GetLogger().GetLevel()
}

func IsTraceLevel() bool {
	return GetLevel() == logrus.TraceLevel
}

func IsDebugLevel() bool {
	return GetLevel() == logrus.DebugLevel
}

func IsInfoLevel() bool {
	return GetLevel() == logrus.InfoLevel
}

func IsWarnLevel() bool {
	return GetLevel() == logrus.WarnLevel
}

func IsErrorLevel() bool {
	return GetLevel() == logrus.ErrorLevel
}

func IsFatalLevel() bool {
	return GetLevel() == logrus.FatalLevel
}

func IsPanicLevel() bool {
	return GetLevel() == logrus.PanicLevel
}

func FormatPanic(value any, stack []byte) string {
	return fmt.Sprintf("panic: %v\n\n%s", value, stack)
}

type loggerBuilder struct {
	logger *Logger
}

func (builder *loggerBuilder) build() *Logger {
	return builder.logger
}

func (builder *loggerBuilder) initialize() *loggerBuilder {
	builder.logger = logrus.New()
	return builder
}

func (builder *loggerBuilder) setSeverity() *loggerBuilder {
	if level, err := logrus.ParseLevel(xbcfg.GetServiceLogLevel()); err != nil {
		panic(err)
	} else {
		builder.logger.SetLevel(level)
	}
	return builder
}

func (builder *loggerBuilder) setFormatter() *loggerBuilder {
	builder.logger.SetFormatter(&jsonifier{})
	return builder
}

type jsonifier struct{}

func (jsonifier *jsonifier) Format(entry *Entry) ([]byte, error) {
	bytes, err := (&jsonBuilder{entry: entry}).
		initialize().
		setLevel().
		setTime().
		setMessage().
		setCaller().
		setFields().
		setBytes().
		build()
	return bytes, err
}

const (
	minCallerFrameSkip = 7
	maxCallerFrameSize = 1 << 5
)

type jsonBuilder struct {
	entry *Entry
	data  Fields
	bytes []byte
	err   error
}

func (builder *jsonBuilder) build() ([]byte, error) {
	return builder.bytes, builder.err
}

func (builder *jsonBuilder) initialize() *jsonBuilder {
	builder.data = make(Fields, 5)
	return builder
}

func (builder *jsonBuilder) setLevel() *jsonBuilder {
	level, _ := builder.entry.Level.MarshalText()
	builder.data["level"] = string(level)
	return builder
}

func (builder *jsonBuilder) setTime() *jsonBuilder {
	builder.data["time"] = builder.entry.Time.Format(time.DateTime + ".000000")
	return builder
}

func (builder *jsonBuilder) setMessage() *jsonBuilder {
	builder.data["message"] = builder.entry.Message
	return builder
}

func (builder *jsonBuilder) setCaller() *jsonBuilder {
	skip := minCallerFrameSkip
	if value, ok := builder.entry.Data[SkipKey]; ok {
		if offset, ok := value.(int); ok {
			skip += offset
		} else {
			panic(fmt.Sprintf("Skip offset `%d` must be of int type.", value))
		}
		delete(builder.entry.Data, SkipKey)
	}
	if caller := getCallerTracer().source(skip); caller != nil {
		callerName := xbslice.Last(strings.Split(caller.Function, "/"))
		builder.data["caller"] = fmt.Sprintf("%s:%d:%s", caller.File, caller.Line, callerName)
	}
	return builder
}

func (builder *jsonBuilder) setFields() *jsonBuilder {
	fields := builder.entry.Data
	if value, ok := fields[ErrorKey]; ok {
		if err, ok := value.(error); ok {
			fields[ErrorKey] = err.Error()
			builder.setErrorFields(fields, err)
		}
	}
	builder.data["fields"] = fields
	return builder
}

func (builder *jsonBuilder) setErrorFields(fields Fields, err error) {
	if uerrs := xberror.Unwrap(err); uerrs != nil {
		for _, uerr := range uerrs {
			builder.setErrorFields(fields, uerr)
		}
	}
	if cerr, ok := xberror.AsCustomError(err); ok {
		maps.Copy(fields, cerr.LogFields())
	}
}

func (builder *jsonBuilder) setBytes() *jsonBuilder {
	bytes, err := xbjson.Marshal(builder.data)
	if err != nil {
		err = xberror.Wrapf("Logger failed to JSON marshal data: `%#v`", []any{builder.data}, err)
	} else {
		bytes = append(bytes, '\n')
	}
	builder.bytes, builder.err = bytes, err
	return builder
}

var mCallerTracer *callerTracer

type callerTracer struct {
	loggerFile  string
	logrusPath  string
	skipFileSet xbctnr.Set[string]
}

func getCallerTracer() *callerTracer {
	if mCallerTracer == nil {
		mCallerTracer = newCallerTracer()
	}
	return mCallerTracer
}

func newCallerTracer() *callerTracer {
	callerTracer := (&callerTracerBuilder{}).
		initialize().
		setLoggerFile().
		setLogrusPath().
		setSkipFileSet().
		build()
	return callerTracer
}

func (tracer *callerTracer) source(skip int) *runtime.Frame {
	pcs := make([]uintptr, maxCallerFrameSize)
	if count := runtime.Callers(skip, pcs); count == 0 {
		return nil
	}
	frames := runtime.CallersFrames(pcs)
	for {
		frame, ok := frames.Next()
		if !ok {
			return nil
		}
		file := frame.File
		if ok := tracer.skipFileSet.Has(file); ok {
			continue
		}
		if strings.Contains(file, tracer.logrusPath) {
			tracer.skipFileSet.Add(file)
			continue
		}
		return &frame
	}
}

type callerTracerBuilder struct {
	callerTracer *callerTracer
}

func (builder *callerTracerBuilder) build() *callerTracer {
	return builder.callerTracer
}

func (builder *callerTracerBuilder) initialize() *callerTracerBuilder {
	builder.callerTracer = &callerTracer{}
	return builder
}

func (builder *callerTracerBuilder) setLoggerFile() *callerTracerBuilder {
	_, file, _, _ := runtime.Caller(0)
	builder.callerTracer.loggerFile = file
	return builder
}

func (builder *callerTracerBuilder) setLogrusPath() *callerTracerBuilder {
	builder.callerTracer.logrusPath = reflect.TypeOf(logrus.Logger{}).PkgPath()
	return builder
}

func (builder *callerTracerBuilder) setSkipFileSet() *callerTracerBuilder {
	set := xbctnr.NewSet[string]()
	set.Add(builder.callerTracer.loggerFile)
	builder.callerTracer.skipFileSet = set
	return builder
}
