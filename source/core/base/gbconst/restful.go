package gbconst

const (
	ContextFlowMap      = "#flow_map"
	FlowKeyFlowID       = "#flow_id"
	FlowKeyFlowTrails   = "#flow_trails"
	FlowKeyFlowError    = "#flow_error"
	FlowKeyFlowOutcome  = "#flow_outcome"
	FlowKeyRequestBody  = "#request_body"
	FlowKeyRequestData  = "#request_data"
	FlowKeyRecordFields = "#record_fields"

	HeaderRealIP             = "X-Real-Ip"
	HeaderForwardedFor       = "X-Forwarded-For"
	HeaderKongRequestID      = "Kong-Request-Id"
	HeaderKongConsumerGroups = "X-Consumer-Groups"

	KongConsumerGroupAnonymous = "anonymous"
	KongConsumerGroupOwner     = "owner"
	KongConsumerGroupClient    = "client"
	KongConsumerGroupService   = "service"
	KongConsumerGroupMonitor   = "monitor"
)
