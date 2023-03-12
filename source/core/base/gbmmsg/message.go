package gbmmsg

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
)

var metaMessageMap = map[string]*MetaMessage{}

type MetaMessage struct {
	metaCode     string
	httpCode     int
	internalText string
	externalCode string
	externalText string
}

func NewMetaMessage(httpCode int, metaCode, externalText, internalText string) *MetaMessage {
	metaMessageMap[metaCode] = (&builder{}).
		initialize().
		setMetaCode(metaCode).
		setHTTPCode(httpCode).
		setInternalText(internalText).
		setExternalCode().
		setExternalText(externalText).
		build()
	return metaMessageMap[metaCode]
}

func (metaMessage *MetaMessage) GetMetaCode() string {
	return metaMessage.metaCode
}

func (metaMessage *MetaMessage) GetHTTPCode() int {
	return metaMessage.httpCode
}

func (metaMessage *MetaMessage) GetInternalText(internalArgs ...interface{}) string {
	return fmt.Sprintf("(%s) %s", metaMessage.metaCode, fmt.Sprintf(metaMessage.internalText, internalArgs...))
}

func (metaMessage *MetaMessage) GetExternalCode() string {
	return metaMessage.externalCode
}

func (metaMessage *MetaMessage) GetExternalText(externalArgs ...interface{}) string {
	return fmt.Sprintf(metaMessage.externalText, externalArgs...)
}

func (metaMessage *MetaMessage) GetExternalStatus() bool {
	return metaMessage.httpCode < http.StatusBadRequest
}

func (metaMessage *MetaMessage) String() string {
	return fmt.Sprintf("<MetaMessage| metaCode: `%s`, httpCode: `%d`>",
		metaMessage.metaCode, metaMessage.httpCode)
}

var metaCodeRegex = regexp.MustCompile(`^[A-Z]{3}[0-9]{3}$`)

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

func (builder *builder) setMetaCode(metaCode string) *builder {
	if _, ok := metaMessageMap[metaCode]; ok {
		panic(fmt.Sprintf("Duplicate meta code `%s` is found.", metaCode))
	}
	if ok := metaCodeRegex.MatchString(metaCode); !ok {
		panic(fmt.Sprintf("Meta code `%s` cannot match regex `%s`.", metaCode, metaCodeRegex.String()))
	}
	builder.metaMessage.metaCode = metaCode
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

func (builder *builder) setInternalText(internalText string) *builder {
	builder.metaMessage.internalText = internalText
	return builder
}

func (builder *builder) setExternalCode() *builder {
	builder.metaMessage.externalCode = fmt.Sprintf("%s-%s", gbcfg.GetServiceCode(), builder.metaMessage.metaCode)
	return builder
}

func (builder *builder) setExternalText(externalText string) *builder {
	builder.metaMessage.externalText = externalText
	return builder
}
