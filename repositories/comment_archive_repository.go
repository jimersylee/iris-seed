package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var CommentArchiveRepository = newCommentArchiveRepository()

func newCommentArchiveRepository() *commentArchiveRepository {
	return &commentArchiveRepository{}
}

type commentArchiveRepository struct {
}

func (this *commentArchiveRepository) Get(db *gorm.DB, id int64) *models.CommentArchive {
	ret := &models.CommentArchive{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *commentArchiveRepository) Take(db *gorm.DB, where ...interface{}) *models.CommentArchive {
	ret := &models.CommentArchive{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *commentArchiveRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.CommentArchive) {
	cnd.Find(db, &list)
	return
}

func (this *commentArchiveRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.CommentArchive) {
	cnd.FindOne(db, &ret)
	return
}

func (this *commentArchiveRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.CommentArchive, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *commentArchiveRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.CommentArchive, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.CommentArchive{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *commentArchiveRepository) Create(db *gorm.DB, t *models.CommentArchive) (err error) {
	err = db.Create(t).Error
	return
}

func (this *commentArchiveRepository) Update(db *gorm.DB, t *models.CommentArchive) (err error) {
	err = db.Save(t).Error
	return
}

func (this *commentArchiveRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.CommentArchive{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *commentArchiveRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.CommentArchive{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *commentArchiveRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.CommentArchive{}, "id = ?", id)
}
