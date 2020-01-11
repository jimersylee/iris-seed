package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var BookRepository = newBookRepository()

func newBookRepository() *bookRepository {
	return &bookRepository{}
}

type bookRepository struct {
}

func (this *bookRepository) Get(db *gorm.DB, id int64) *models.Book {
	ret := &models.Book{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *bookRepository) Take(db *gorm.DB, where ...interface{}) *models.Book {
	ret := &models.Book{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *bookRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Book) {
	cnd.Find(db, &list)
	return
}

func (this *bookRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.Book) {
	cnd.FindOne(db, &ret)
	return
}

func (this *bookRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Book, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *bookRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Book, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Book{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *bookRepository) Create(db *gorm.DB, t *models.Book) (err error) {
	err = db.Create(t).Error
	return
}

func (this *bookRepository) Update(db *gorm.DB, t *models.Book) (err error) {
	err = db.Save(t).Error
	return
}

func (this *bookRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Book{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *bookRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Book{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *bookRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Book{}, "id = ?", id)
}
