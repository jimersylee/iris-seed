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


func JsonSuccess() *WebApiResult {
	return &WebApiResult{
		ErrorCode: 0,
		Data:      nil,
		Success:   true,
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
