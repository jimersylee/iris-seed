package commons

import (
	"encoding/json"
)

var (
	ErrorCodeSystem         = NewError(1000, "系统异常")
	ErrorCodeNotLogin       = NewError(1, "请先登录")
	ErrorCodeParse          = NewError(2, "解析错误")
	ErrorCodeNotFound       = NewError(3, "未找到")
	ErrorCodeRegisterFailed = NewError(4, "用户注册失败")
)

func NewError(code int, text string) *ErrorCode {
	return &ErrorCode{code, text, nil, false, nil}
}
func NewErrorData(code int, text string, errorData interface{}) *ErrorCode {
	return &ErrorCode{code, text, nil, false, errorData}
}

func FromError(err error) *ErrorCode {
	if err == nil {
		return nil
	}
	return &ErrorCode{0, err.Error(), nil, false, nil}
}

type ErrorCode struct {
	Code      int
	Message   string
	Data      interface{}
	Success   bool
	ErrorData interface{}
}

func (e *ErrorCode) Error() string {
	marshal, _ := json.Marshal(e)
	return string(marshal)
}
