package api

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
)

type ArticleController struct {
	Ctx iris.Context
}

func (this *ArticleController) GetBy(id int64) *response.WebApiRes {
	t := services.ArticleService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(t)
}

func (this *ArticleController) AnyList() *response.WebApiRes {
	list, paging := services.ArticleService.FindPageByParams(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return response.JsonData(&commons.PageResult{Results: list, Page: paging})
}

func (this *ArticleController) PostCreate() *response.WebApiRes {
	t := &models.Article{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.ArticleService.Create(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}

func (this *ArticleController) PostUpdate() *response.WebApiRes {
	id, err := commons.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	t := services.ArticleService.Get(id)
	if t == nil {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	err = this.Ctx.ReadForm(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}

	err = services.ArticleService.Update(t)
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	return response.JsonData(t)
}
func (this *ArticleController) GetTagBy(name string) *response.WebApiRes {

	page := commons.NewQueryParams(this.Ctx).PageByReq().Desc("id").Paging.Page
	limit := commons.NewQueryParams(this.Ctx).PageByReq().Desc("id").Paging.Limit
	data, paging := services.ArticleService.FindByTag(name, page, limit)
	return response.JsonData(&commons.PageResult{Results: data, Page: paging})

}

//
// GetArchiveBy
// @Description:
// @receiver this
// @param yearmonth:如:2017-01
// @return *response.WebApiRes
//
func (this *ArticleController) GetArchiveBy(yearmonth string) *response.WebApiRes {
	data, page := services.ArticleService.FindArticlesByYearMonth(commons.NewQueryParams(this.Ctx).PageByReq().Desc("id"), yearmonth)
	return response.JsonData(&commons.PageResult{Results: data, Page: page})

}

//
// GetTags
// @Description: 获取所有标签
// @receiver this
// @return *response.WebApiRes
//
func (this *ArticleController) GetTags() *response.WebApiRes {
	tags := services.ArticleService.FindTags()
	return response.JsonData(tags)
}

//
// GetGuess
// @Description: 猜你喜欢
// @receiver this
// @return *response.WebApiRes
//
func (this *ArticleController) GetGuess() *response.WebApiRes {
	likeArticles := services.ArticleService.GuessYouLike()
	return response.JsonData(likeArticles)
}

//
// GetHot
// @Description: 热门文章
// @receiver this
// @return *response.WebApiRes
//
func (this *ArticleController) GetHot() *response.WebApiRes {
	article := services.ArticleService.FindHot()
	return response.JsonData(article)
}

//
// GetCategory
// @Description: 分类目录
// @receiver this
// @return *response.WebApiRes
//
func (this *ArticleController) GetCategory() *response.WebApiRes {
	category := services.ArticleService.FindCategory()
	return response.JsonData(category)
}
