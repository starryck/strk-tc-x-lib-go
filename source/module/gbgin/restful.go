package gbgin

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbconst"
	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbmtmsg"
	"github.com/forbot161602/pbc-golang-lib/source/core/toolkit/gbjson"
	"github.com/forbot161602/pbc-golang-lib/source/core/toolkit/gbslice"
	"github.com/forbot161602/pbc-golang-lib/source/core/utility/gberror"
	"github.com/forbot161602/pbc-golang-lib/source/core/utility/gblogger"
)

type RESTFlow struct {
	BaseFlow
	context *Context
}

func (flow *RESTFlow) Initiate(context *Context) {
	flow.context = context
	if flowMap, ok := context.Get(gbconst.ContextFlowMap); ok {
		storage, _ := flowMap.(*sync.Map)
		flow.SetStorage(storage)
	} else {
		flow.BaseFlow.Initiate()
		context.Set(gbconst.ContextFlowMap, flow.GetStorage())
	}
	return
}

func (flow *RESTFlow) Inherit(fore Flow) {
	panic("REST flow doesn't support inheritance.")
}

func (flow *RESTFlow) SetError(err error) {
	flow.context.Abort()
	flow.context.Error(err)
	flow.BaseFlow.SetError(err)
	return
}

func (flow *RESTFlow) SetNotFoundError() {
	flow.SetError(gberror.Validation(gbmtmsg.WMV400, nil))
	return
}

func (flow *RESTFlow) GetContext() *Context {
	context := flow.context
	return context
}

func (flow *RESTFlow) GetRequest() *http.Request {
	request := flow.context.Request
	return request
}

func (flow *RESTFlow) GetRequestIP() string {
	ip := ""
	if ip = strings.Split(flow.GetHeader(gbconst.HeaderForwardedFor), ",")[0]; ip != "" {
		return ip
	}
	if ip = flow.GetHeader(gbconst.HeaderRealIP); ip != "" {
		return ip
	}
	if ip = flow.context.ClientIP(); ip != "" {
		return ip
	}
	return ip
}

func (flow *RESTFlow) GetRequestURI() string {
	uri := flow.context.Request.RequestURI
	return uri
}

func (flow *RESTFlow) GetMethod() string {
	method := flow.context.Request.Method
	return method
}

func (flow *RESTFlow) IsGetMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodGet)
}

func (flow *RESTFlow) IsHeadMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodHead)
}

func (flow *RESTFlow) IsPutMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodPut)
}

func (flow *RESTFlow) IsPostMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodPost)
}

func (flow *RESTFlow) IsPatchMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodPatch)
}

func (flow *RESTFlow) IsDeleteMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodDelete)
}

func (flow *RESTFlow) IsTraceMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodTrace)
}

func (flow *RESTFlow) IsOptionsMethod() bool {
	method := flow.GetMethod()
	return strings.EqualFold(method, http.MethodOptions)
}

func (flow *RESTFlow) GetParam(key string) string {
	value := flow.context.Param(key)
	return value
}

func (flow *RESTFlow) GetQuery(key string) string {
	value := flow.context.Query(key)
	return value
}

func (flow *RESTFlow) GetQueryMap(key string) map[string]string {
	value := flow.context.QueryMap(key)
	return value
}

func (flow *RESTFlow) GetQueryArray(key string) []string {
	value := flow.context.QueryArray(key)
	return value
}

func (flow *RESTFlow) GetQueryValues() url.Values {
	value := flow.context.Request.URL.Query()
	return value
}

func (flow *RESTFlow) GetQueryFallback(key, fallback string) string {
	value := flow.context.DefaultQuery(key, fallback)
	return value
}

func (flow *RESTFlow) GetHeader(key string) string {
	value := flow.context.GetHeader(key)
	return value
}

func (flow *RESTFlow) GetHeaderValues() http.Header {
	value := flow.context.Request.Header
	return value
}

func (flow *RESTFlow) GetBody() io.ReadCloser {
	body := flow.context.Request.Body
	return body
}

func (flow *RESTFlow) BindBody(value any) {
	body := flow.GetBody()
	if body == nil {
		flow.SetError(gberror.Validation(gbmtmsg.WMV420, &gberror.Options{
			LogFields: gblogger.Fields{
				"requestBody": body,
			},
		}))
		return
	}
	data, _ := io.ReadAll(body)
	if err := gbjson.Unmarshal(data, value); err != nil {
		flow.SetError(gberror.Validation(gbmtmsg.WMV421, &gberror.Options{
			LogFields: gblogger.Fields{
				"requestData": string(data),
			},
		}))
		return
	}
	flow.Expose(gbconst.FlowKeyRequestBody, value)
	flow.Expose(gbconst.FlowKeyRequestData, data)
	return
}

func (flow *RESTFlow) ContainBody() bool {
	ok := flow.Contain(gbconst.FlowKeyRequestBody)
	return ok
}

func (flow *RESTFlow) RequireBody() any {
	body := flow.Require(gbconst.FlowKeyRequestBody)
	return body
}

func (flow *RESTFlow) ContainData() bool {
	ok := flow.Contain(gbconst.FlowKeyRequestData)
	return ok
}

func (flow *RESTFlow) RequireData() []byte {
	data := flow.RequireBytes(gbconst.FlowKeyRequestData)
	return data
}

func (flow *RESTFlow) GetWriter() ResponseWriter {
	writer := flow.context.Writer
	return writer
}

func (flow *RESTFlow) SetHeader(key, value string) {
	flow.context.Header(key, value)
	return
}

func (flow *RESTFlow) RespondFile(path string) {
	flow.context.File(path)
	return
}

func (flow *RESTFlow) RespondJSON(message *MetaMessage, data any, options *JSONResponseOptions) {
	response := NewJSONResponse(message, data, options)
	flow.context.JSON(response.Code, response)
	return
}

type KongFlow struct {
	RESTFlow
}

func (flow *KongFlow) GetRequestID() string {
	id := flow.GetHeader(gbconst.HeaderKongRequestID)
	return id
}

func (flow *KongFlow) GetConsumerGroups() []string {
	groups := strings.Split(flow.GetHeader(gbconst.HeaderKongConsumerGroups), ",")
	return groups
}

func (flow *KongFlow) IsAnonymousRequest() bool {
	groups := flow.GetConsumerGroups()
	return gbslice.Contain(groups, gbconst.KongConsumerGroupAnonymous)
}

func (flow *KongFlow) IsOwnerRequest() bool {
	groups := flow.GetConsumerGroups()
	return gbslice.Contain(groups, gbconst.KongConsumerGroupOwner)
}

func (flow *KongFlow) IsClientRequest() bool {
	groups := flow.GetConsumerGroups()
	return gbslice.Contain(groups, gbconst.KongConsumerGroupClient)
}

func (flow *KongFlow) IsServiceRequest() bool {
	groups := flow.GetConsumerGroups()
	return gbslice.Contain(groups, gbconst.KongConsumerGroupService)
}

func (flow *KongFlow) IsMonitorRequest() bool {
	groups := flow.GetConsumerGroups()
	return gbslice.Contain(groups, gbconst.KongConsumerGroupMonitor)
}

func (flow *KongFlow) IsInternalRequest() bool {
	isValid := flow.IsServiceRequest() || flow.IsMonitorRequest()
	return isValid
}

func (flow *KongFlow) IsExternalRequest() bool {
	isValid := flow.IsOwnerRequest() || flow.IsClientRequest()
	return isValid
}

func (flow *KongFlow) IsAuthenticatedRequest() bool {
	isValid := flow.IsInternalRequest() || flow.IsExternalRequest()
	return isValid
}
