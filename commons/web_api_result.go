package commons

type WebApiResult struct {
	ErrorCode int         `json:"errorCode"`
	ErrorData interface{} `json:"errorData"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Success   bool        `json:"success"`
}

func Json(code int, message string, data interface{}, success bool) *WebApiResult {
	return &WebApiResult{
		ErrorCode: code,
		Message:   message,
		Data:      data,
		Success:   success,
	}
}

func JsonData(data interface{}) *WebApiResult {
	return &WebApiResult{
		ErrorCode: 0,
		Data:      data,
		Success:   true,
	}
}

func JsonPageData(results interface{}, page *Paging) *WebApiResult {
	return JsonData(&PageResult{
		Results: results,
		Page:    page,
	})
}

func JsonCursorData(results interface{}, cursor string) *WebApiResult {
	return JsonData(&CursorResult{
		Results: results,
		Cursor:  cursor,
	})
}

func JsonSuccess() *WebApiResult {
	return &WebApiResult{
		ErrorCode: 0,
		Data:      nil,
		Success:   true,
	}
}

func JsonError(err *CodeError) *WebApiResult {
	return &WebApiResult{
		ErrorCode: err.Code,
		Message:   err.Message,
		Data:      err.Data,
		Success:   false,
	}
}

func JsonErrorMsg(message string) *WebApiResult {
	return &WebApiResult{
		ErrorCode: 0,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorCode(code int, message string) *WebApiResult {
	return &WebApiResult{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorData(code int, message string, data interface{}) *WebApiResult {
	return &WebApiResult{
		ErrorCode: code,
		Message:   message,
		Data:      data,
		Success:   false,
	}
}

type RspBuilder struct {
	Data map[string]interface{}
}

func NewEmptyRspBuilder() *RspBuilder {
	return &RspBuilder{Data: make(map[string]interface{})}
}

func NewRspBuilder(obj interface{}) *RspBuilder {
	return NewRspBuilderExcludes(obj)
}

func NewRspBuilderExcludes(obj interface{}, excludes ...string) *RspBuilder {
	return &RspBuilder{Data: StructToMap(obj, excludes...)}
}

func (this *RspBuilder) Put(key string, value interface{}) *RspBuilder {
	this.Data[key] = value
	return this
}

func (this *RspBuilder) Build() map[string]interface{} {
	return this.Data
}

func (this *RspBuilder) JsonResult() *WebApiResult {
	return JsonData(this.Data)
}
