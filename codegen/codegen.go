package main

import (
	"github.com/jimersylee/go-steam-proxy/commons/code_generator"
	"github.com/jimersylee/go-steam-proxy/models"
)

func main() {
	code_generator.Generate("./", "github.com/jimersylee/go-steam-proxy", code_generator.GetGenerateStruct(&models.Ip{}))
}
