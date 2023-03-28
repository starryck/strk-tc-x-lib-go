package gbmtmsg

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbcfg"
)

// Code: {Severity Code (1)}{Project Code (1)}{Service Code (1)}{Sequence Number (3)}
// Severity code: F(fetal), E(error), W(warning), I(info), ...
// Project code: A(project1), B(project2), ..., M(module1), N(module2), ...
// Service code: E(entry), V(view), S(service), M(model), K(kernel), ...
// Sequence number: 000, 001, 002, 003, 004, 005, ...
var (
	// RESTful view
	WMV400 = NewMetaMessage(http.StatusBadRequest,
		"WMV400", "RESTful view: Bad request.",
		"Bad request.")
	WMV403 = NewMetaMessage(http.StatusForbidden,
		"WMV403", "RESTful view: Forbidden.",
		"Forbidden.")
	WMV404 = NewMetaMessage(http.StatusNotFound,
		"WMV404", "RESTful view: Not found.",
		"Not found.")
	EMV500 = NewMetaMessage(http.StatusInternalServerError,
		"EMV500", "RESTful view: Internal server error.",
		"Internal server error.")
	WMV420 = NewMetaMessage(http.StatusBadRequest,
		"WMV420", "RESTful view: Invalid parameter.",
		"Request body must be not nil.")
	WMV421 = NewMetaMessage(http.StatusBadRequest,
		"WMV421", "RESTful view: Invalid parameter.",
		"Request body must be bound correctly.")
)

var metaMessageMap = map[string]*MetaMessage{}

func NewMetaMessage(httpCode int, code, outText, logText string) *MetaMessage {
	metaMessageMap[code] = (&builder{}).
		initialize().
		setCode(code).
		setHTTPCode(httpCode).
		setLogText(logText).
		setOutCode().
		setOutText(outText).
		build()
	return metaMessageMap[code]
}

type MetaMessage struct {
	code     string
	httpCode int
	logText  string
	outCode  string
	outText  string
}

func (metaMessage *MetaMessage) GetCode() string {
	return metaMessage.code
}

func (metaMessage *MetaMessage) GetHTTPCode() int {
	return metaMessage.httpCode
}

func (metaMessage *MetaMessage) GetLogText(logArgs ...any) string {
	return fmt.Sprintf("(%s) %s", metaMessage.code, fmt.Sprintf(metaMessage.logText, logArgs...))
}

func (metaMessage *MetaMessage) GetOutCode() string {
	return metaMessage.outCode
}

func (metaMessage *MetaMessage) GetOutText(outArgs ...any) string {
	return fmt.Sprintf(metaMessage.outText, outArgs...)
}

func (metaMessage *MetaMessage) GetOutStatus() bool {
	return metaMessage.httpCode < http.StatusBadRequest
}

func (metaMessage *MetaMessage) String() string {
	return fmt.Sprintf("<MetaMessage| code: `%s`, httpCode: `%d`>",
		metaMessage.code, metaMessage.httpCode)
}

var metaMessageCodeRegex = regexp.MustCompile(`^[A-Z]{3}[0-9]{3}$`)

type builder struct {
	metaMessage *MetaMessage
}

func (builder *builder) build() *MetaMessage {
	return builder.metaMessage
}

func (builder *builder) initialize() *builder {
	builder.metaMessage = &MetaMessage{}
	return builder
}

func (builder *builder) setCode(code string) *builder {
	if _, ok := metaMessageMap[code]; ok {
		panic(fmt.Sprintf("Duplicate meta message code `%s` is found.", code))
	}
	if ok := metaMessageCodeRegex.MatchString(code); !ok {
		panic(fmt.Sprintf("Meta message code `%s` cannot match regex `%s`.", code, metaMessageCodeRegex.String()))
	}
	builder.metaMessage.code = code
	return builder
}

func (builder *builder) setHTTPCode(httpCode int) *builder {
	minCode, maxCode := http.StatusContinue, http.StatusNetworkAuthenticationRequired
	if httpCode < minCode || httpCode > maxCode {
		panic(fmt.Sprintf("HTTP code `%d` must be between `%d` and `%d`.", httpCode, minCode, maxCode))
	}
	builder.metaMessage.httpCode = httpCode
	return builder
}

func (builder *builder) setLogText(logText string) *builder {
	builder.metaMessage.logText = logText
	return builder
}

func (builder *builder) setOutCode() *builder {
	builder.metaMessage.outCode = fmt.Sprintf("%s-%s", gbcfg.GetServiceCode(), builder.metaMessage.code)
	return builder
}

func (builder *builder) setOutText(outText string) *builder {
	builder.metaMessage.outText = outText
	return builder
}
