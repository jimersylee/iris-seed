package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var LoginHistoryService = newLoginHistoryService()

func newLoginHistoryService() *loginHistoryService {
	return &loginHistoryService{}
}

type loginHistoryService struct {
}

func (this *loginHistoryService) Get(id int64) *models.LoginHistory {
	return repositories.LoginHistoryRepository.Get(db.GetDB(), id)
}

func (this *loginHistoryService) Take(where ...interface{}) *models.LoginHistory {
	return repositories.LoginHistoryRepository.Take(db.GetDB(), where...)
}

func (this *loginHistoryService) Find(cnd *commons.SqlCnd) []models.LoginHistory {
	return repositories.LoginHistoryRepository.Find(db.GetDB(), cnd)
}

func (this *loginHistoryService) FindOne(cnd *commons.SqlCnd) *models.LoginHistory {
	return repositories.LoginHistoryRepository.FindOne(db.GetDB(), cnd)
}

func (this *loginHistoryService) FindPageByParams(params *commons.QueryParams) (list []models.LoginHistory, paging *commons.Paging) {
	return repositories.LoginHistoryRepository.FindPageByParams(db.GetDB(), params)
}

func (this *loginHistoryService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.LoginHistory, paging *commons.Paging) {
	return repositories.LoginHistoryRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *loginHistoryService) Create(t *models.LoginHistory) error {
	return repositories.LoginHistoryRepository.Create(db.GetDB(), t)
}

func (this *loginHistoryService) Update(t *models.LoginHistory) error {
	return repositories.LoginHistoryRepository.Update(db.GetDB(), t)
}

func (this *loginHistoryService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.LoginHistoryRepository.Updates(db.GetDB(), id, columns)
}

func (this *loginHistoryService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.LoginHistoryRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *loginHistoryService) Delete(id int64) {
	repositories.LoginHistoryRepository.Delete(db.GetDB(), id)
}
