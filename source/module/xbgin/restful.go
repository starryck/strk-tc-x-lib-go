package xbgin

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/forbot161602/x-lib-go/source/core/base/xbconst"
	"github.com/forbot161602/x-lib-go/source/core/base/xbmtmsg"
	"github.com/forbot161602/x-lib-go/source/core/toolkit/xbjson"
	"github.com/forbot161602/x-lib-go/source/core/toolkit/xbslice"
	"github.com/forbot161602/x-lib-go/source/core/utility/xberror"
	"github.com/forbot161602/x-lib-go/source/core/utility/xblogger"
)

type RESTFlow struct {
	BaseFlow
	context *Context
}

func (flow *RESTFlow) Initiate(context *Context) {
	flow.context = context
	if flowMap, ok := context.Get(xbconst.ContextFlowMap); ok {
		storage, _ := flowMap.(*sync.Map)
		flow.SetStorage(storage)
	} else {
		flow.BaseFlow.Initiate()
		context.Set(xbconst.ContextFlowMap, flow.GetStorage())
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
	flow.SetError(xberror.Validation(xbmtmsg.WMV404, nil))
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
	if ip = strings.Split(flow.GetHeader(xbconst.HeaderForwardedFor), ",")[0]; ip != "" {
		return ip
	}
	if ip = flow.GetHeader(xbconst.HeaderRealIP); ip != "" {
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
		flow.SetError(xberror.Validation(xbmtmsg.WMV420, &xberror.Options{
			LogFields: xblogger.Fields{
				"requestBody": body,
			},
		}))
		return
	}
	data, _ := io.ReadAll(body)
	if err := xbjson.Unmarshal(data, value); err != nil {
		flow.SetError(xberror.Validation(xbmtmsg.WMV421, &xberror.Options{
			LogFields: xblogger.Fields{
				"requestData": string(data),
			},
		}))
		return
	}
	flow.Expose(xbconst.FlowKeyRequestBody, value)
	flow.Expose(xbconst.FlowKeyRequestData, data)
	return
}

func (flow *RESTFlow) ContainBody() bool {
	ok := flow.Contain(xbconst.FlowKeyRequestBody)
	return ok
}

func (flow *RESTFlow) RequireBody() any {
	body := flow.Require(xbconst.FlowKeyRequestBody)
	return body
}

func (flow *RESTFlow) ContainData() bool {
	ok := flow.Contain(xbconst.FlowKeyRequestData)
	return ok
}

func (flow *RESTFlow) RequireData() []byte {
	data := flow.RequireBytes(xbconst.FlowKeyRequestData)
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
	id := flow.GetHeader(xbconst.HeaderKongRequestID)
	return id
}

func (flow *KongFlow) GetConsumerGroups() []string {
	groups := strings.Split(flow.GetHeader(xbconst.HeaderKongConsumerGroups), ",")
	return groups
}

func (flow *KongFlow) IsAnonymousRequest() bool {
	groups := flow.GetConsumerGroups()
	return xbslice.Contain(groups, xbconst.KongConsumerGroupAnonymous)
}

func (flow *KongFlow) IsOwnerRequest() bool {
	groups := flow.GetConsumerGroups()
	return xbslice.Contain(groups, xbconst.KongConsumerGroupOwner)
}

func (flow *KongFlow) IsClientRequest() bool {
	groups := flow.GetConsumerGroups()
	return xbslice.Contain(groups, xbconst.KongConsumerGroupClient)
}

func (flow *KongFlow) IsServiceRequest() bool {
	groups := flow.GetConsumerGroups()
	return xbslice.Contain(groups, xbconst.KongConsumerGroupService)
}

func (flow *KongFlow) IsMonitorRequest() bool {
	groups := flow.GetConsumerGroups()
	return xbslice.Contain(groups, xbconst.KongConsumerGroupMonitor)
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
