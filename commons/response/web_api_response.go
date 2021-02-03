package response

import "github.com/jimersylee/iris-seed/commons"

type WebApiRes struct {
	ErrorCode int         `json:"errorCode"`
	ErrorData interface{} `json:"errorData"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Success   bool        `json:"success"`
}

func Json(code int, message string, data interface{}, success bool) *WebApiRes {
	return &WebApiRes{
		ErrorCode: code,
		Message:   message,
		Data:      data,
		Success:   success,
	}
}

func JsonData(data interface{}) *WebApiRes {
	return &WebApiRes{
		ErrorCode: 0,
		Data:      data,
		Success:   true,
		Message:   "success",
	}
}

func JsonErrorCode(errorCode *commons.ErrorCode) *WebApiRes {
	return &WebApiRes{
		ErrorCode: errorCode.Code,
		ErrorData: nil,
		Message:   errorCode.Message,
		Data:      nil,
		Success:   false,
	}
}

func JsonSuccess() *WebApiRes {
	return &WebApiRes{
		ErrorCode: 0,
		Data:      nil,
		Success:   true,
	}
}

func JsonErrorMsg(message string) *WebApiRes {
	return &WebApiRes{
		ErrorCode: 0,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonError(code int, message string) *WebApiRes {
	return &WebApiRes{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorData(code int, message string, errorData interface{}) *WebApiRes {
	return &WebApiRes{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
		Success:   false,
		ErrorData: errorData,
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

func (this *RspBuilder) JsonResult() *WebApiRes {
	return JsonData(this.Data)
}
