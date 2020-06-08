package api

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/services"
	"github.com/jimersylee/iris-seed/services/cache"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

type IpController struct {
	Ctx iris.Context
}

func (this *IpController) GetBy(id int64) *response.WebApiRes {
	t := services.IpService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(t)
}

func (this *IpController) AnyList() *response.WebApiRes {
	list, paging := services.IpService.FindPageByParams(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return response.JsonData(&commons.PageResult{Results: list, Page: paging})
}

func (this *IpController) PostCreate() *response.WebApiRes {
	t := &models.Ip{}
	err := this.Ctx.ReadJSON(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.IpService.Create(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

func (this *IpController) PostUpdate() *response.WebApiRes {
	id, err := commons.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	t := services.IpService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.IpService.Update(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

//统计数据
func (this *IpController) GetStat() *response.WebApiRes {
	logrus.Info(cache.ProxyCache.GetStatistic())
	return response.JsonSuccess()
}
