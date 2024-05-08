package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"path"
	"strconv"
	"strings"
)

/*
*
GET /index.html HTTP/1.1
Host: localhost:4221
User-Agent: curl/7.64.1
*/
func parseMessage(buf []byte) bool {
	msg := string(buf)
	startLine := strings.Split(msg, "\r\n")[0]
	path := strings.Split(startLine, " ")[1]
	method := strings.Split(startLine, " ")[0]
	protocol := strings.Split(startLine, " ")[2]
	fmt.Println("method,path,protocol:", method, path, protocol)
	if path == "/" {
		return true
	} else {
		return false
	}
}

type Handler func(req *Request, res *Response)

type Server struct {
	handlers map[string]Handler
}

func (s *Server) SetHandler(pattern string, handler Handler) {
	s.handlers[pattern] = handler
}

func (s *Server) Start(port int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	log.Printf("server start at port %d", port)
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

const HeaderConnection = "Connection"

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	res := NewResponse(StatusOK)
	res.Headers[HeaderConnection] = "close"
	res.Headers[HeaderContentLength] = "0"
	req, err := readRequest(conn)
	if err != nil {
		res.Status = StatusBadRequest
	} else {
		handler := s.getHandler(req.Path)
		if handler != nil {
			handler(req, res)
		} else {
			res.Status = StatusNotFound
		}
	}
	sendResponse(conn, res)
}

func sendResponse(conn net.Conn, res *Response) {
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("HTTP/1.1 %s\r\n", res.Status))
	for key, value := range res.Headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n")
	msg.WriteString(res.Content)
	conn.Write([]byte(msg.String()))
}

func (s *Server) getHandler(requestPath string) Handler {
	// todo: implement path pattern-matching
	handler, ok := s.handlers[requestPath]
	if ok {
		return handler
	}
	if requestPath != "/" {
		requestPath = path.Dir(requestPath)
		handler, ok := s.handlers[requestPath+"/**"]
		if ok {
			return handler
		}
	}
	return nil
}

func readRequest(conn net.Conn) (*Request, error) {
	req := Request{Headers: make(map[string]string)}

	reader := bufio.NewReader(conn)
	requestLine, err := readLine(reader)

	if err != nil {
		return nil, err
	}
	reqParts := strings.Fields(string(requestLine))
	if len(reqParts) != 3 {
		return nil, fmt.Errorf("invalid req line: %s", string(requestLine))
	}
	req.Method = reqParts[0]
	req.Path = reqParts[1]
	req.HttpVersion = reqParts[2]

	for {
		headerLine, err := readLine(reader)
		if err != nil {
			return nil, err
		} else if len(headerLine) == 0 {
			break
		}
		headerParts := strings.SplitN(string(headerLine), ": ", 2)
		if len(headerParts) != 2 {
			continue
		}
		req.Headers[headerParts[0]] = headerParts[1]
	}

	if req.Headers[HeaderContentLength] == "" {
		return &req, nil
	}

	contentLength, err := strconv.Atoi(req.Headers[HeaderContentLength])
	if err != nil {
		return nil, err
	}
	body := strings.Builder{}
	for i := 0; i < contentLength; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		body.WriteByte(b)
	}
	req.Body = body.String()
	return &req, nil
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		next, prefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line = append(line, next...)
		if !prefix {
			break
		}
	}

	return line, nil
}

func NewServer() *Server {
	return &Server{
		handlers: make(map[string]Handler),
	}
}
