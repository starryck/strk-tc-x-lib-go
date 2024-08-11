package xbdata

type PaginationResult struct {
	PageIndex   int `json:"pageIndex"`
	PageSize    int `json:"pageSize"`
	PageCount   int `json:"pageCount"`
	RecordCount int `json:"recordCount"`
}

type JSONResponseMeta struct {
	JSONResponseBaseData
}

type JSONResponsePageMeta struct {
	JSONResponseBaseData
	JSONResponsePageData
}

type JSONResponseBaseData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type JSONResponsePageData = PaginationResult
