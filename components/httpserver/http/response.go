package http

import "strconv"

const StatusOK = "200 OK"
const StatusNotFound = "404 Not Found"
const StatusBadRequest = "400 Bad Request"
const StatusCreated = "201 Created"
const StatusInternalServerErr = "500 Internal Server Error"

type Response struct {
	Status  string
	Headers map[string]string
	Content string
}

func NewResponse(status string) *Response {
	return &Response{
		Status:  status,
		Headers: make(map[string]string),
	}
}

func (r *Response) SetContent(contentType, content string) {
	r.Content = content
	r.Headers[HeaderContentLength] = strconv.Itoa(len(content))
	r.Headers[HeaderContentType] = contentType
}
