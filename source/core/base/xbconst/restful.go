package xbconst

const (
	ContextFlowMap        = "#flow_map"
	FlowKeyFlowID         = "#flow_id"
	FlowKeyFlowTrails     = "#flow_trails"
	FlowKeyFlowError      = "#flow_error"
	FlowKeyFlowOutcome    = "#flow_outcome"
	FlowKeyRequestParams  = "#request_params"
	FlowKeyRequestQueries = "#request_queries"
	FlowKeyRequestHeaders = "#request_headers"
	FlowKeyRequestBody    = "#request_body"
	FlowKeyRequestData    = "#request_data"
	FlowKeyRecordFields   = "#record_fields"

	HeaderAuthorization = "Authorization"
	HeaderRealIP        = "X-Real-Ip"
	HeaderForwardedFor  = "X-Forwarded-For"

	HeaderKongRequestID        = "Kong-Request-Id"
	HeaderKongConsumerCustomID = "X-Consumer-Custom-Id"
	HeaderKongConsumerGroups   = "X-Consumer-Groups"
	KongConsumerGroupAnonymous = "anonymous"
	KongConsumerGroupUser      = "owner"
	KongConsumerGroupClient    = "client"
	KongConsumerGroupService   = "service"
	KongConsumerGroupMonitor   = "monitor"

	HeaderAPISIXRequestID        = "X-Request-Id"
	HeaderAPISIXConsumerName     = "X-Consumer-Name"
	HeaderAPISIXConsumerGroupID  = "X-Consumer-Group-Id"
	APISIXConsumerGroupIDUser    = "dft_user"
	APISIXConsumerGroupIDClient  = "dft_client"
	APISIXConsumerGroupIDService = "dft_service"
	APISIXConsumerGroupIDMonitor = "dft_monitor"
)
