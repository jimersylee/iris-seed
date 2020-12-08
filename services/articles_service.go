package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var ArticlesService = newArticlesService()

func newArticlesService() *articlesService {
	return &articlesService{}
}

type articlesService struct {
}

func (this *articlesService) Get(id int64) *models.Articles {
	return repositories.ArticlesRepository.Get(db.GetDB(), id)
}

func (this *articlesService) Take(where ...interface{}) *models.Articles {
	return repositories.ArticlesRepository.Take(db.GetDB(), where...)
}

func (this *articlesService) Find(cnd *commons.SqlCnd) []models.Articles {
	return repositories.ArticlesRepository.Find(db.GetDB(), cnd)
}

func (this *articlesService) FindOne(cnd *commons.SqlCnd) *models.Articles {
	return repositories.ArticlesRepository.FindOne(db.GetDB(), cnd)
}

func (this *articlesService) FindPageByParams(params *commons.QueryParams) (list []models.Articles, paging *commons.Paging) {
	return repositories.ArticlesRepository.FindPageByParams(db.GetDB(), params)
}

func (this *articlesService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Articles, paging *commons.Paging) {
	return repositories.ArticlesRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *articlesService) Create(t *models.Articles) error {
	return repositories.ArticlesRepository.Create(db.GetDB(), t)
}

func (this *articlesService) Update(t *models.Articles) error {
	return repositories.ArticlesRepository.Update(db.GetDB(), t)
}

func (this *articlesService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ArticlesRepository.Updates(db.GetDB(), id, columns)
}

func (this *articlesService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ArticlesRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *articlesService) Delete(id int64) {
	repositories.ArticlesRepository.Delete(db.GetDB(), id)
}
