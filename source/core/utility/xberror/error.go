package xberror

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/starryck/strk-tc-x-lib-go/source/core/base/xbmtmsg"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func AsError(err error) (target *Error, ok bool) {
	return target, As(err, &target)
}

func AsWrapError(err error) (target *WrapError, ok bool) {
	return target, As(err, &target)
}

func AsNestedError(err error) (target NestedError, ok bool) {
	return target, As(err, &target)
}

func AsCustomError(err error) (target CustomError, ok bool) {
	return target, As(err, &target)
}

func AsInternalError(err error) (target *InternalError, ok bool) {
	return target, As(err, &target)
}

func AsValidationError(err error) (target *ValidationError, ok bool) {
	return target, As(err, &target)
}

func AsUnexpectedError(err error) (target *UnexpectedError, ok bool) {
	return target, As(err, &target)
}

func Unwrap(err error) []error {
	if uerr := errors.Unwrap(err); uerr != nil {
		return []error{uerr}
	}
	if nerr, ok := err.(NestedError); ok {
		return nerr.Unwrap()
	}
	return nil
}

func Aggravate(err error) error {
	if verr, ok := AsValidationError(err); ok {
		err = Unexpected(verr.Message(), verr.Options(), verr.Unwrap()...)
	}
	return err
}

type Error struct {
	message string
}

func New(message string) *Error {
	return &Error{
		message: message,
	}
}

func Newf(message string, args []any) *Error {
	return New(fmt.Sprintf(message, args...))
}

func (err *Error) Error() string {
	return fmt.Sprintf("<Error| %s>", err.message)
}

type WrapError struct {
	message string
	errs    []error
}

func Wrap(message string, errs ...error) *WrapError {
	return &WrapError{
		message: message,
		errs:    errs,
	}
}

func Wrapf(message string, args []any, errs ...error) *WrapError {
	return Wrap(fmt.Sprintf(message, args...), errs...)
}

func (err *WrapError) Error() string {
	return fmt.Sprintf("<WrapError| %s> is caused from: %v", err.message, err.errs)
}

func (err *WrapError) Unwrap() []error {
	return err.errs
}

type (
	Message   = xbmtmsg.MetaMessage
	Options   = internalErrorOptions
	LogFields = logrus.Fields
)

type NestedError interface {
	Error() string
	Unwrap() []error
}

type CustomError interface {
	NestedError
	Message() *Message
	Options() *Options
	OutText() string
	OutArgs() []any
	LogText() string
	LogArgs() []any
	LogFields() LogFields
}

type InternalError struct {
	message   *Message
	options   *Options
	errs      []error
	outArgs   []any
	logArgs   []any
	logFields LogFields
}

func Internal(message *Message, options *Options, errs ...error) *InternalError {
	err := (&internalErrorBuilder{options: options}).
		initialize().
		setMessage(message).
		setOptions().
		setErrors(errs...).
		setOutArgs().
		setLogArgs().
		setLogFields().
		build()
	return err
}

func (err *InternalError) Error() string {
	return fmt.Sprintf("<InternalError| %s> is caused from: %v", err.LogText(), err.errs)
}

func (err *InternalError) Unwrap() []error {
	return err.errs
}

func (err *InternalError) Message() *Message {
	return err.message
}

func (err *InternalError) Options() *Options {
	return err.options
}

func (err *InternalError) OutText() string {
	return err.message.GetOutText(err.outArgs...)
}

func (err *InternalError) OutArgs() []any {
	return err.outArgs
}

func (err *InternalError) LogText() string {
	return err.message.GetLogText(err.logArgs...)
}

func (err *InternalError) LogArgs() []any {
	return err.logArgs
}

func (err *InternalError) LogFields() LogFields {
	return err.logFields
}

type internalErrorBuilder struct {
	err     *InternalError
	options *Options
}

type internalErrorOptions struct {
	OutArgs   []any
	LogArgs   []any
	LogFields LogFields
}

func (builder *internalErrorBuilder) build() *InternalError {
	return builder.err
}

func (builder *internalErrorBuilder) initialize() *internalErrorBuilder {
	builder.err = &InternalError{}
	if builder.options == nil {
		builder.options = &Options{}
	}
	return builder
}

func (builder *internalErrorBuilder) setMessage(message *Message) *internalErrorBuilder {
	builder.err.message = message
	return builder
}

func (builder *internalErrorBuilder) setOptions() *internalErrorBuilder {
	builder.err.options = builder.options
	return builder
}

func (builder *internalErrorBuilder) setErrors(errs ...error) *internalErrorBuilder {
	builder.err.errs = errs
	return builder
}

func (builder *internalErrorBuilder) setOutArgs() *internalErrorBuilder {
	builder.err.outArgs = builder.options.OutArgs
	return builder
}

func (builder *internalErrorBuilder) setLogArgs() *internalErrorBuilder {
	builder.err.logArgs = builder.options.LogArgs
	return builder
}

func (builder *internalErrorBuilder) setLogFields() *internalErrorBuilder {
	builder.err.logFields = builder.options.LogFields
	return builder
}

type ValidationError struct {
	*InternalError
}

func Validation(message *Message, options *Options, errs ...error) *ValidationError {
	return &ValidationError{
		InternalError: Internal(message, options, errs...),
	}
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("<ValidationError| %s> is caused from: %v", err.LogText(), err.errs)
}

type UnexpectedError struct {
	*InternalError
}

func Unexpected(message *Message, options *Options, errs ...error) *UnexpectedError {
	return &UnexpectedError{
		InternalError: Internal(message, options, errs...),
	}
}

func (err *UnexpectedError) Error() string {
	return fmt.Sprintf("<UnexpectedError| %s> is caused from: %v", err.LogText(), err.errs)
}
