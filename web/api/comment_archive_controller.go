package api

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
)

type CommentArchiveController struct {
	Ctx iris.Context
}

func (this *CommentArchiveController) GetBy(id int64) *response.WebApiRes {
	t := services.CommentArchiveService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(t)
}

func (this *CommentArchiveController) AnyList() *response.WebApiRes {
	list, paging := services.CommentArchiveService.FindPageByParams(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return response.JsonData(&commons.PageResult{Results: list, Page: paging})
}

func (this *CommentArchiveController) PostCreate() *response.WebApiRes {
	t := &models.CommentArchive{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.CommentArchiveService.Create(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

func (this *CommentArchiveController) PostUpdate() *response.WebApiRes {
	id, err := commons.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	t := services.CommentArchiveService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.CommentArchiveService.Update(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}
