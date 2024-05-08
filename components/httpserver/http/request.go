package http

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

type Request struct {
	Path        string
	Method      string
	HttpVersion string
	Headers     map[string]string
	Body        string
}
