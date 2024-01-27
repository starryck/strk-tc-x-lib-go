package xbmodel

type JSONResponseMeta struct {
	JSONResponseBaseData
}

type JSONResponsePageMeta struct {
	JSONResponseBaseData
	JSONResponsePageData
}

type JSONResponseBaseData struct {
	Code    string `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type JSONResponsePageData struct {
	PageIndex   int `json:"pageIndex"`
	PageSize    int `json:"pageSize"`
	PageConut   int `json:"pageCount"`
	RecordCount int `json:"recordCount"`
}
