package xbgin

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/forbot161602/x-lib-go/source/core/base/xbconst"
	"github.com/forbot161602/x-lib-go/source/core/base/xbmtmsg"
	"github.com/forbot161602/x-lib-go/source/core/toolkit/xbslice"
	"github.com/forbot161602/x-lib-go/source/core/utility/xberror"
	"github.com/forbot161602/x-lib-go/source/core/utility/xblogger"
	"github.com/forbot161602/x-lib-go/source/core/utility/xbwatch"
	"github.com/forbot161602/x-lib-go/source/utility/xbspvs"
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
	flow.SetTimer()
	flow.SetBodies()
	flow.NextFlow()
	flow.SetFields()
	flow.SetResult()
}

type RecordMiddlewareFlow struct {
	MiddlewareFlow
	timer  *xbwatch.Timer
	bodies []byte
	fields xblogger.Fields
}

func (flow *RecordMiddlewareFlow) SetTimer() {
	flow.timer = xbwatch.NewTimer()
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
	return
}

func (flow *RecordMiddlewareFlow) SetFields() {
	fields := xblogger.Fields{
		"requestIP":      flow.makeRequestIP(),
		"requestURI":     flow.makeRequestURI(),
		"requestMethod":  flow.makeRequestMethod(),
		"requestHandler": flow.makeRequestHandler(),
		"requestContent": flow.makeRequestContent(),
		"responseTime":   flow.makeResponseTime(),
		"responseSize":   flow.makeResponseSize(),
		"responseStatus": flow.makeResponseStatus(),
	}
	flow.fields = fields
	flow.Expose(xbconst.FlowKeyRecordFields, fields)
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
	time := float64(flow.timer.Stamp().ElapsedTimeMs()) / 1000
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
