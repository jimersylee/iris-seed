package api

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
)

type AdminController struct {
	Ctx iris.Context
}

func (this *AdminController) GetBy(id int64) *response.WebApiRes {
	t := services.AdminService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(t)
}

func (this *AdminController) AnyList() *response.WebApiRes {
	list, paging := services.AdminService.FindPageByParams(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return response.JsonData(&commons.PageResult{Results: list, Page: paging})
}

func (this *AdminController) PostCreate() *response.WebApiRes {
	t := &models.Admin{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.AdminService.Create(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

func (this *AdminController) PostUpdate() *response.WebApiRes {
	id, err := commons.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	t := services.AdminService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.AdminService.Update(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}
