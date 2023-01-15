package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var CommentRepository = newCommentRepository()

func newCommentRepository() *commentRepository {
	return &commentRepository{}
}

type commentRepository struct {
}

func (this *commentRepository) Get(db *gorm.DB, id int64) *models.Comment {
	ret := &models.Comment{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *commentRepository) Take(db *gorm.DB, where ...interface{}) *models.Comment {
	ret := &models.Comment{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *commentRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Comment) {
	cnd.Find(db, &list)
	return
}

func (this *commentRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.Comment) {
	cnd.FindOne(db, &ret)
	return
}

func (this *commentRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Comment, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *commentRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Comment, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Comment{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *commentRepository) Create(db *gorm.DB, t *models.Comment) (err error) {
	err = db.Create(t).Error
	return
}

func (this *commentRepository) Update(db *gorm.DB, t *models.Comment) (err error) {
	err = db.Save(t).Error
	return
}

func (this *commentRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Comment{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *commentRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *commentRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Comment{}, "id = ?", id)
}
