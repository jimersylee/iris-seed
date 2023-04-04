package entities

// BizException
// @Description: 业务异常,比如用户不存在,用户已经存在等,都可以panic
type BizException struct {
	Status int `json:"status"`

	Code string `json:"code"`

	Message string `json:"message"`

	Trace string `json:"trace"`
}

func (e *BizException) Error() string {

	return e.Trace

}
