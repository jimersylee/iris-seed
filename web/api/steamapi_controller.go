package api

import (
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/kataras/iris"
)

type SteamapiController struct {
	Ctx iris.Context
}

func (c *SteamapiController) Post() *response.WebApiRes {
	return response.JsonSuccess()

}
