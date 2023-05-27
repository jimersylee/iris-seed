package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var RoleRepository = newRoleRepository()

func newRoleRepository() *roleRepository {
	return &roleRepository{}
}

type roleRepository struct {
}

func (this *roleRepository) Get(db *gorm.DB, id int64) *models.Role {
	ret := &models.Role{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *roleRepository) Take(db *gorm.DB, where ...interface{}) *models.Role {
	ret := &models.Role{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *roleRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Role) {
	cnd.Find(db, &list)
	return
}

func (this *roleRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.Role) {
	cnd.FindOne(db, &ret)
	return
}

func (this *roleRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Role, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *roleRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Role, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Role{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *roleRepository) Create(db *gorm.DB, t *models.Role) (err error) {
	err = db.Create(t).Error
	return
}

func (this *roleRepository) Update(db *gorm.DB, t *models.Role) (err error) {
	err = db.Save(t).Error
	return
}

func (this *roleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Role{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *roleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Role{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *roleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Role{}, "id = ?", id)
}
