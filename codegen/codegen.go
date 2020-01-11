package main

import (
	"github.com/jimersylee/iris-seed/commons/codegenerator"
	"github.com/jimersylee/iris-seed/models"
)

func main() {
	codegenerator.Generate("./", "github.com/jimersylee/iris-seed", codegenerator.GetGenerateStruct(&models.Book{}))
}
