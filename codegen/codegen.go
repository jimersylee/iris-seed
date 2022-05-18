package main

import (
	"github.com/jimersylee/iris-seed/commons/code_generator"
	"github.com/jimersylee/iris-seed/models"
)

/**
 * 这是 iris-seed 代码生成器
 */
func main() {
	code_generator.Generate("./", "github.com/jimersylee/iris-seed", code_generator.GetGenerateStruct(&models.Article{}))
}
