package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var RoleService = newRoleService()

func newRoleService() *roleService {
	return &roleService{}
}

type roleService struct {
}

func (this *roleService) Get(id int64) *models.Role {
	return repositories.RoleRepository.Get(db.GetDB(), id)
}

func (this *roleService) Take(where ...interface{}) *models.Role {
	return repositories.RoleRepository.Take(db.GetDB(), where...)
}

func (this *roleService) Find(cnd *commons.SqlCnd) []models.Role {
	return repositories.RoleRepository.Find(db.GetDB(), cnd)
}

func (this *roleService) FindOne(cnd *commons.SqlCnd) *models.Role {
	return repositories.RoleRepository.FindOne(db.GetDB(), cnd)
}

func (this *roleService) FindPageByParams(params *commons.QueryParams) (list []models.Role, paging *commons.Paging) {
	return repositories.RoleRepository.FindPageByParams(db.GetDB(), params)
}

func (this *roleService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Role, paging *commons.Paging) {
	return repositories.RoleRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *roleService) Create(t *models.Role) error {
	return repositories.RoleRepository.Create(db.GetDB(), t)
}

func (this *roleService) Update(t *models.Role) error {
	return repositories.RoleRepository.Update(db.GetDB(), t)
}

func (this *roleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.RoleRepository.Updates(db.GetDB(), id, columns)
}

func (this *roleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.RoleRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *roleService) Delete(id int64) {
	repositories.RoleRepository.Delete(db.GetDB(), id)
}
