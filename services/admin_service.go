package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var AdminService = newAdminService()

func newAdminService() *adminService {
	return &adminService{}
}

type adminService struct {
}

func (this *adminService) Get(id int64) *models.Admin {
	return repositories.AdminRepository.Get(db.GetDB(), id)
}

func (this *adminService) Take(where ...interface{}) *models.Admin {
	return repositories.AdminRepository.Take(db.GetDB(), where...)
}

func (this *adminService) Find(cnd *commons.SqlCnd) []models.Admin {
	return repositories.AdminRepository.Find(db.GetDB(), cnd)
}

func (this *adminService) FindOne(cnd *commons.SqlCnd) *models.Admin {
	return repositories.AdminRepository.FindOne(db.GetDB(), cnd)
}

func (this *adminService) FindPageByParams(params *commons.QueryParams) (list []models.Admin, paging *commons.Paging) {
	return repositories.AdminRepository.FindPageByParams(db.GetDB(), params)
}

func (this *adminService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Admin, paging *commons.Paging) {
	return repositories.AdminRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *adminService) Create(t *models.Admin) error {
	return repositories.AdminRepository.Create(db.GetDB(), t)
}

func (this *adminService) Update(t *models.Admin) error {
	return repositories.AdminRepository.Update(db.GetDB(), t)
}

func (this *adminService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.AdminRepository.Updates(db.GetDB(), id, columns)
}

func (this *adminService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.AdminRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *adminService) Delete(id int64) {
	repositories.AdminRepository.Delete(db.GetDB(), id)
}
