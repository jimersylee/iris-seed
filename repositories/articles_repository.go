package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var ArticlesRepository = newArticlesRepository()

func newArticlesRepository() *articlesRepository {
	return &articlesRepository{}
}

type articlesRepository struct {
}

func (this *articlesRepository) Get(db *gorm.DB, id int64) *models.Articles {
	ret := &models.Articles{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articlesRepository) Take(db *gorm.DB, where ...interface{}) *models.Articles {
	ret := &models.Articles{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articlesRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Articles) {
	cnd.Find(db, &list)
	return
}

func (this *articlesRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.Articles) {
	cnd.FindOne(db, &ret)
	return
}

func (this *articlesRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Articles, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *articlesRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Articles, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Articles{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *articlesRepository) Create(db *gorm.DB, t *models.Articles) (err error) {
	err = db.Create(t).Error
	return
}

func (this *articlesRepository) Update(db *gorm.DB, t *models.Articles) (err error) {
	err = db.Save(t).Error
	return
}

func (this *articlesRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Articles{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *articlesRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Articles{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *articlesRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Articles{}, "id = ?", id)
}
