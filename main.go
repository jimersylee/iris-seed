package main

import (
	_ "net/http/pprof"

	"github.com/jimersylee/iris-seed/app"
)

// @author:jimersylee@gmail.com
func main() {
	app.RunApp()

}
