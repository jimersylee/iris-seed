package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var NoteRepository = newNoteRepository()

func newNoteRepository() *noteRepository {
	return &noteRepository{}
}

type noteRepository struct {
}

func (this *noteRepository) Get(db *gorm.DB, id int64) *models.Note {
	ret := &models.Note{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *noteRepository) Take(db *gorm.DB, where ...interface{}) *models.Note {
	ret := &models.Note{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *noteRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Note) {
	cnd.Find(db, &list)
	return
}

func (this *noteRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.Note) {
	cnd.FindOne(db, &ret)
	return
}

func (this *noteRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Note, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *noteRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Note, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Note{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *noteRepository) Create(db *gorm.DB, t *models.Note) (err error) {
	err = db.Create(t).Error
	return
}

func (this *noteRepository) Update(db *gorm.DB, t *models.Note) (err error) {
	err = db.Save(t).Error
	return
}

func (this *noteRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Note{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *noteRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Note{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *noteRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Note{}, "id = ?", id)
}
