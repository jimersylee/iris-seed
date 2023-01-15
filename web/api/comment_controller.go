package api

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
)

type CommentController struct {
	Ctx iris.Context
}

func (this *CommentController) GetBy(id int64) *response.WebApiRes {
	t := services.CommentService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(t)
}

func (this *CommentController) AnyList() *response.WebApiRes {
	list, paging := services.CommentService.FindPageByParams(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return response.JsonData(&commons.PageResult{Results: list, Page: paging})
}

func (this *CommentController) PostCreate() *response.WebApiRes {
	t := &models.Comment{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.CommentService.Create(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

func (this *CommentController) PostUpdate() *response.WebApiRes {
	id, err := commons.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	t := services.CommentService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.CommentService.Update(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}
