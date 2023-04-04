package api

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
)

type PaymentHistoryController struct {
	Ctx iris.Context
}

func (this *PaymentHistoryController) GetBy(id int64) *response.WebApiRes {
	t := services.PaymentHistoryService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(t)
}

func (this *PaymentHistoryController) AnyList() *response.WebApiRes {
	list, paging := services.PaymentHistoryService.FindPageByParams(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return response.JsonData(&commons.PageResult{Results: list, Page: paging})
}

func (this *PaymentHistoryController) PostCreate() *response.WebApiRes {
	t := &models.PaymentHistory{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.PaymentHistoryService.Create(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

func (this *PaymentHistoryController) PostUpdate() *response.WebApiRes {
	id, err := commons.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	t := services.PaymentHistoryService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.PaymentHistoryService.Update(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}
