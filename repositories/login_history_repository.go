package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var LoginHistoryRepository = newLoginHistoryRepository()

func newLoginHistoryRepository() *loginHistoryRepository {
	return &loginHistoryRepository{}
}

type loginHistoryRepository struct {
}

func (this *loginHistoryRepository) Get(db *gorm.DB, id int64) *models.LoginHistory {
	ret := &models.LoginHistory{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *loginHistoryRepository) Take(db *gorm.DB, where ...interface{}) *models.LoginHistory {
	ret := &models.LoginHistory{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *loginHistoryRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.LoginHistory) {
	cnd.Find(db, &list)
	return
}

func (this *loginHistoryRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.LoginHistory) {
	cnd.FindOne(db, &ret)
	return
}

func (this *loginHistoryRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.LoginHistory, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *loginHistoryRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.LoginHistory, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.LoginHistory{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *loginHistoryRepository) Create(db *gorm.DB, t *models.LoginHistory) (err error) {
	err = db.Create(t).Error
	return
}

func (this *loginHistoryRepository) Update(db *gorm.DB, t *models.LoginHistory) (err error) {
	err = db.Save(t).Error
	return
}

func (this *loginHistoryRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.LoginHistory{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *loginHistoryRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.LoginHistory{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *loginHistoryRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.LoginHistory{}, "id = ?", id)
}
