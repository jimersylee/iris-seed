package main

import (
	"github.com/jimersylee/iris-seed/commons/code_generator"
	"github.com/jimersylee/iris-seed/models"
)

/**
 * easy to generator some code
 */
func main() {
	code_generator.Generate("./", "github.com/jimersylee/iris-seed", code_generator.GetGenerateStruct(&models.Admin{}))
}
