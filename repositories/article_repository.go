package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/entities"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var ArticleRepository = newArticleRepository()

func newArticleRepository() *articleRepository {
	return &articleRepository{}
}

type articleRepository struct {
}

func (this *articleRepository) Get(db *gorm.DB, id int64) *models.Article {
	ret := &models.Article{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articleRepository) Take(db *gorm.DB, where ...interface{}) *models.Article {
	ret := &models.Article{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articleRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Article) {
	cnd.Find(db, &list)
	return
}

func (this *articleRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.Article) {
	err := cnd.FindOne(db, &ret)
	if err != nil {
		return nil
	}
	return
}

func (this *articleRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Article, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *articleRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Article, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Article{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *articleRepository) Create(db *gorm.DB, t *models.Article) (err error) {
	err = db.Create(t).Error
	return
}

func (this *articleRepository) Update(db *gorm.DB, t *models.Article) (err error) {
	err = db.Save(t).Error
	return
}

func (this *articleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Article{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *articleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *articleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Article{}, "id = ?", id)
}

func (this *articleRepository) FindArticlesByTagId(db *gorm.DB, id int64) (list []models.Article, paging *commons.Paging) {

	var count int
	db.Preload("Tags").Model(&models.Article{}).Where("id in (select article_id from article_tag where tag_id = ?)", id).Count(&count)
	db.Preload("Tags").Model(&models.Article{}).Where("id in (select article_id from article_tag where tag_id = ?)", id).Find(&list)

	paging = &commons.Paging{
		Page:  1,
		Limit: 1,
		Total: count,
	}
	return

}
func (this *articleRepository) FindArticlesByYeahMonth(db *gorm.DB, startTime int, endTime int, page int, limit int) (list []models.Article, paging *commons.Paging) {
	var count int

	db.Preload("Tags").Model(&models.Article{}).Where("create_at >= ? and create_at < ?", startTime, endTime).Offset((page - 1) * limit).Limit(limit).Find(&list)
	db.Preload("Tags").Model(&models.Article{}).Where("create_at >= ? and create_at < ?", startTime, endTime).Count(&count)

	paging = &commons.Paging{
		Page:  page,
		Limit: limit,
		Total: count,
	}
	return
}

func (this *articleRepository) FindArticlesByTag(db *gorm.DB, tag string, page int, limit int) (list []models.Article, paging *commons.Paging) {

	return
}

func (this *articleRepository) FindCategory(db *gorm.DB) (list []entities.TagDTO) {
	db.Raw("select category as name, count(category) as num from article group by category;").Find(&list)
	return
}
