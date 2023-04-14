package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var AdminRepository = newAdminRepository()

func newAdminRepository() *adminRepository {
	return &adminRepository{}
}

type adminRepository struct {
}

func (this *adminRepository) Get(db *gorm.DB, id int64) *models.Admin {
	ret := &models.Admin{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *adminRepository) Take(db *gorm.DB, where ...interface{}) *models.Admin {
	ret := &models.Admin{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *adminRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Admin) {
	cnd.Find(db, &list)
	return
}

func (this *adminRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.Admin) {
	cnd.FindOne(db, &ret)
	return
}

func (this *adminRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Admin, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *adminRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Admin, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Admin{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *adminRepository) Create(db *gorm.DB, t *models.Admin) (err error) {
	err = db.Create(t).Error
	return
}

func (this *adminRepository) Update(db *gorm.DB, t *models.Admin) (err error) {
	err = db.Save(t).Error
	return
}

func (this *adminRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Admin{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *adminRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Admin{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *adminRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Admin{}, "id = ?", id)
}
