package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/entities"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
	"time"
)

var ArticleService = newArticleService()

func newArticleService() *articleService {
	return &articleService{}
}

type articleService struct {
}

func (this *articleService) Get(id int64) *models.Article {
	return repositories.ArticleRepository.Get(db.GetDB(), id)
}

func (this *articleService) Take(where ...interface{}) *models.Article {
	return repositories.ArticleRepository.Take(db.GetDB(), where...)
}

func (this *articleService) Find(cnd *commons.SqlCnd) []models.Article {
	return repositories.ArticleRepository.Find(db.GetDB(), cnd)
}

func (this *articleService) FindOne(cnd *commons.SqlCnd) *models.Article {
	return repositories.ArticleRepository.FindOne(db.GetDB(), cnd)
}

func (this *articleService) FindPageByParams(params *commons.QueryParams) (list []models.Article, paging *commons.Paging) {
	return repositories.ArticleRepository.FindPageByParams(db.GetDB(), params)
}

func (this *articleService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Article, paging *commons.Paging) {
	return repositories.ArticleRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *articleService) Create(t *models.Article) error {
	return repositories.ArticleRepository.Create(db.GetDB(), t)
}

func (this *articleService) Update(t *models.Article) error {
	return repositories.ArticleRepository.Update(db.GetDB(), t)
}

func (this *articleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ArticleRepository.Updates(db.GetDB(), id, columns)
}

func (this *articleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ArticleRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *articleService) Delete(id int64) {
	repositories.ArticleRepository.Delete(db.GetDB(), id)
}

func (this *articleService) FindArticlesByYearMonth(queryParam *commons.QueryParams, yearMonth string) (list []models.Article, paging *commons.Paging) {
	//时间转时间戳
	t, err := time.ParseInLocation("2006-01", yearMonth, time.Local)
	if err != nil {
		panic(commons.ErrorCodeParse)
	}
	//获取当月第一天时间戳
	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
	//获取当月最后一天时间戳
	endTime := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Unix()
	articles, paging := repositories.ArticleRepository.FindArticlesByYeahMonth(db.GetDB(), int(startTime), int(endTime), queryParam.Paging.Page, queryParam.Paging.Limit)
	if articles == nil {
		return nil, paging
	}
	return articles, paging
}

func (this *articleService) FindTags() (list []entities.TagDTO) {
	db.GetDB().Raw("select a.tag as name,count(a.tag) as num from article, json_table(tags_string,'$[*]' columns (tag varchar(40) PATH '$')) a group by a.tag;").Find(&list)
	return list
}

func (this *articleService) GuessYouLike() (list []models.Article) {
	RecommendService.GetRecommendArticles(1, 10)
	return
}

func (this *articleService) FindByTag(name string, page int, limit int) (list []models.Article, paging *commons.Paging) {
	db.GetDB().Raw("select * from article where json_contains(tags_string->'$','\"" + name + "\"')").Limit(limit).Offset((page - 1) * limit).Order("id desc").Find(&list)
	count := 0
	db.GetDB().Model(models.Article{}).Where("json_contains(tags_string->'$','\"" + name + "\"')").Count(&count)

	paging = &commons.Paging{
		Page:  page,
		Limit: limit,
		Total: count,
	}
	return
}

//
// FindHot
// @Description: 查询热门文章
// @receiver this
// @return interface{}
//
func (this *articleService) FindHot() (list []models.Article) {
	cnd := commons.NewSqlCnd()
	cnd.Orders = append(cnd.Orders, commons.OrderByCol{Column: "view_times", Asc: false})
	cnd.Limit(10)
	return repositories.ArticleRepository.Find(db.GetDB(), cnd)
}

//
// FindCategory
// @Description: 查询分类文章
// @receiver this
// @return list
//
func (this *articleService) FindCategory() (list []entities.TagDTO) {
	return repositories.ArticleRepository.FindCategory(db.GetDB())
}
