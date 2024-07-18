package xbgin

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/starryck/x-lib-go/source/core/base/xbconst"
	"github.com/starryck/x-lib-go/source/core/base/xbmtmsg"
	"github.com/starryck/x-lib-go/source/core/toolkit/xbslice"
	"github.com/starryck/x-lib-go/source/core/utility/xberror"
	"github.com/starryck/x-lib-go/source/core/utility/xblogger"
	"github.com/starryck/x-lib-go/source/core/utility/xbwatch"
	"github.com/starryck/x-lib-go/source/utility/xbspvs"
)

type MiddlewareFlow struct {
	RESTFlow
}

func (flow *MiddlewareFlow) NextFlow() {
	flow.GetContext().Next()
	return
}

func GraceMiddleware(ctx *Context) {
	flow := &GraceMiddlewareFlow{}
	flow.Initiate(ctx)
	flow.NextFlow()
}

type GraceMiddlewareFlow struct {
	MiddlewareFlow
}

func (flow *GraceMiddlewareFlow) NextFlow() {
	xbspvs.WithWaitGroup(func(args ...any) {
		flow.MiddlewareFlow.NextFlow()
	})
	return
}

const (
	MaxRequestBodyReadSize   = 1 << 12
	MaxRequestBodyRecordSize = 1 << 16
)

func RecordMiddleware(ctx *Context) {
	flow := &RecordMiddlewareFlow{}
	flow.Initiate(ctx)
	flow.SetBodies()
	flow.NextFlow()
	flow.SetFields()
	flow.SetResult()
}

type RecordMiddlewareFlow struct {
	MiddlewareFlow
	watch  *xbwatch.Watch
	fields xblogger.Fields
	bodies []byte
}

func (flow *RecordMiddlewareFlow) Initiate(ctx *Context) {
	flow.MiddlewareFlow.Initiate(ctx)
	flow.watch = xbwatch.NewWatch()
	flow.fields = xblogger.Fields{}
	flow.Expose(xbconst.FlowKeyRecordFields, flow.fields)
	return
}

func (flow *RecordMiddlewareFlow) SetBodies() {
	request := flow.GetRequest()
	if request.Body == nil {
		return
	}

	buffer := &bytes.Buffer{}
	bodies := make([]byte, MaxRequestBodyRecordSize)
	if length, _ := request.Body.Read(bodies); length > 0 {
		buffer.Write(bodies[:length])
		flow.bodies = buffer.Bytes()
	}
	for {
		bodies := make([]byte, MaxRequestBodyReadSize)
		if length, _ := request.Body.Read(bodies); length > 0 {
			buffer.Write(bodies[:length])
		} else {
			break
		}
	}
	request.Body = io.NopCloser(buffer)
	flow.Expose(xbconst.FlowKeyRequestData, buffer.Bytes())
	return
}

func (flow *RecordMiddlewareFlow) SetFields() {
	fields := flow.fields
	fields["requestIP"] = flow.makeRequestIP()
	fields["requestURI"] = flow.makeRequestURI()
	fields["requestMethod"] = flow.makeRequestMethod()
	fields["requestHandler"] = flow.makeRequestHandler()
	fields["requestContent"] = flow.makeRequestContent()
	fields["responseTime"] = flow.makeResponseTime()
	fields["responseSize"] = flow.makeResponseSize()
	fields["responseStatus"] = flow.makeResponseStatus()
	return
}

func (flow *RecordMiddlewareFlow) makeRequestIP() string {
	ip := flow.GetRequestIP()
	return ip
}

func (flow *RecordMiddlewareFlow) makeRequestURI() string {
	uri := flow.GetRequestURI()
	return uri
}

func (flow *RecordMiddlewareFlow) makeRequestMethod() string {
	method := flow.GetMethod()
	return method
}

func (flow *RecordMiddlewareFlow) makeRequestHandler() string {
	handler := xbslice.Last(strings.Split(flow.GetContext().HandlerName(), "/"))
	return handler
}

func (flow *RecordMiddlewareFlow) makeRequestContent() string {
	content := string(flow.bodies)
	return content
}

func (flow *RecordMiddlewareFlow) makeResponseTime() float64 {
	time := float64(flow.watch.Stamp().ElapsedTimeMs()) / 1000
	return time
}

func (flow *RecordMiddlewareFlow) makeResponseSize() int {
	size := flow.GetContext().Writer.Size()
	return size
}

func (flow *RecordMiddlewareFlow) makeResponseStatus() int {
	status := flow.GetContext().Writer.Status()
	return status
}

func (flow *RecordMiddlewareFlow) SetResult() {
	flow.GetLogger().WithFields(flow.fields).Info("Log access message.")
	if err := flow.GetError(); err != nil {
		if flow.makeResponseStatus() >= http.StatusInternalServerError {
			flow.GetLogger().WithError(err).Error("Log error message.")
		} else {
			flow.GetLogger().WithError(err).Warning("Log warning message.")
		}
	}
	return
}

func ResponseMiddleware(ctx *Context) {
	flow := &ResponseMiddlewareFlow{}
	flow.Initiate(ctx)
	flow.NextFlow()
	if flow.HasError() {
		flow.SetResult()
	}
}

type ResponseMiddlewareFlow struct {
	MiddlewareFlow
}

func (flow *ResponseMiddlewareFlow) SetResult() {
	err := flow.GetError()
	if cerr, ok := xberror.AsCustomError(err); ok {
		flow.RespondJSON(cerr.Message(), nil, &JSONResponseOptions{
			MetaArgs: cerr.OutArgs(),
		})
		return
	}
	flow.RespondJSON(xbmtmsg.EMV500, nil, nil)
	return
}
