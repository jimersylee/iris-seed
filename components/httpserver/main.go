package main

import (
	"flag"
	"github.com/jimersylee/iris-seed/components/httpserver/http"
	"os"
	"path"
	"strings"
)

func main() {
	var directory string
	flag.StringVar(&directory, "directory", "files", "path to the file directory")
	flag.Parse()
	server := http.NewServer()
	server.SetHandler("/", func(req *http.Request, res *http.Response) {
	})
	server.SetHandler("/echo/**", func(req *http.Request, res *http.Response) {
		index := strings.Index(req.Path, "/echo/")
		if index > -1 {
			res.SetContent("text/plain", req.Path[index+6:])
		}
	})
	server.SetHandler("/user-agent", func(req *http.Request, res *http.Response) {
		res.SetContent("text/plain", req.Headers[http.HeaderUserAgent])
	})

	server.SetHandler("/files/**", func(req *http.Request, res *http.Response) {
		if req.Method == http.MethodGet {
			filePath := path.Join(directory, req.Path[len("/files/"):])
			file, err := os.ReadFile(filePath)
			if err != nil {
				res.Status = http.StatusNotFound
			} else {
				res.SetContent("application/octet-stream", string(file))
			}
		} else if req.Method == http.MethodPost {
			filePath := path.Join(directory, req.Path[len("/files/"):])
			err := os.WriteFile(filePath, []byte(req.Body), 0644)
			if err != nil {
				res.Status = http.StatusBadRequest
			} else {
				res.Status = http.StatusCreated
			}
		}

	})

	err := server.Start(4221)
	if err != nil {
		panic(err)
	}

}
