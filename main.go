package main

import (
	"fmt"
	"github.com/jimersylee/iris-seed/app"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// @author:jimersylee@gmail.com
func main() {
	go func() {
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()
	app.RunApp()

}
