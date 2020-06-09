package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var IpRepository = newIpRepository()

func newIpRepository() *ipRepository {
	return &ipRepository{}
}

type ipRepository struct {
}

func (this *ipRepository) Get(db *gorm.DB, id int64) *models.Ip {
	ret := &models.Ip{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *ipRepository) Take(db *gorm.DB, where ...interface{}) *models.Ip {
	ret := &models.Ip{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *ipRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Ip) {
	cnd.Find(db, &list)
	return
}

func (this *ipRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) *models.Ip {
	ret := &models.Ip{}
	err := cnd.FindOne(db, &ret)
	if err != nil {
		return nil
	}
	return ret
}

func (this *ipRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Ip, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *ipRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Ip, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Ip{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *ipRepository) Create(db *gorm.DB, t *models.Ip) (err error) {
	err = db.Create(t).Error
	return
}

func (this *ipRepository) Update(db *gorm.DB, t *models.Ip) (err error) {
	err = db.Save(t).Error
	return
}

func (this *ipRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Ip{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *ipRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Ip{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *ipRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Ip{}, "id = ?", id)
}
