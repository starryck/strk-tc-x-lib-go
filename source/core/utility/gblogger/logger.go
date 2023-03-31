package gblogger

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbcfg"
	"github.com/forbot161602/pbc-golang-lib/source/core/toolkit/gbjson"
	"github.com/forbot161602/pbc-golang-lib/source/core/toolkit/gbslice"
	"github.com/forbot161602/pbc-golang-lib/source/core/utility/gberror"
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
	logger := (&builder{}).
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

type builder struct {
	logger *Logger
}

func (builder *builder) build() *Logger {
	return builder.logger
}

func (builder *builder) initialize() *builder {
	builder.logger = logrus.New()
	return builder
}

func (builder *builder) setSeverity() *builder {
	if level, err := logrus.ParseLevel(gbcfg.GetServiceLogLevel()); err != nil {
		panic(err)
	} else {
		builder.logger.SetLevel(level)
	}
	return builder
}

func (builder *builder) setFormatter() *builder {
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
	MinCallerFrameSkip = 7
	MaxCallerFrameSize = 1 << 5
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
	skip := MinCallerFrameSkip
	if value, ok := builder.entry.Data[SkipKey]; ok {
		if offset, ok := value.(int); ok {
			skip += offset
		} else {
			panic(fmt.Sprintf("Skip offset `%d` must be of int type.", value))
		}
		delete(builder.entry.Data, SkipKey)
	}
	if caller := getCallerTracer().source(skip); caller != nil {
		callerFunc := gbslice.Last(strings.Split(caller.Function, "/"))
		builder.data["caller"] = fmt.Sprintf("%s:%d:%s", caller.File, caller.Line, callerFunc)
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
	if uerrs := gberror.Unwrap(err); uerrs != nil {
		for _, uerr := range uerrs {
			builder.setErrorFields(fields, uerr)
		}
	}
	if cerr, ok := gberror.AsCustomError(err); ok {
		for key, value := range cerr.LogFields() {
			fields[key] = value
		}
	}
}

func (builder *jsonBuilder) setBytes() *jsonBuilder {
	bytes, err := gbjson.Marshal(builder.data)
	if err != nil {
		err = gberror.Wrapf("Logger failed to JSON marshal data: `%#v`", []any{builder.data}, err)
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
	skipFileSet map[string]bool
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
	pcs := make([]uintptr, MaxCallerFrameSize)
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
		if _, ok := tracer.skipFileSet[file]; ok {
			continue
		}
		if strings.Contains(file, tracer.logrusPath) {
			tracer.skipFileSet[file] = true
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
	fileSet := make(map[string]bool, MaxCallerFrameSize)
	fileSet[builder.callerTracer.loggerFile] = true
	builder.callerTracer.skipFileSet = fileSet
	return builder
}
