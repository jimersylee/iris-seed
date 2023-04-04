package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var PaymentHistoryService = newPaymentHistoryService()

func newPaymentHistoryService() *paymentHistoryService {
	return &paymentHistoryService{}
}

type paymentHistoryService struct {
}

func (this *paymentHistoryService) Get(id int64) *models.PaymentHistory {
	return repositories.PaymentHistoryRepository.Get(db.GetDB(), id)
}

func (this *paymentHistoryService) Take(where ...interface{}) *models.PaymentHistory {
	return repositories.PaymentHistoryRepository.Take(db.GetDB(), where...)
}

func (this *paymentHistoryService) Find(cnd *commons.SqlCnd) []models.PaymentHistory {
	return repositories.PaymentHistoryRepository.Find(db.GetDB(), cnd)
}

func (this *paymentHistoryService) FindOne(cnd *commons.SqlCnd) *models.PaymentHistory {
	return repositories.PaymentHistoryRepository.FindOne(db.GetDB(), cnd)
}

func (this *paymentHistoryService) FindPageByParams(params *commons.QueryParams) (list []models.PaymentHistory, paging *commons.Paging) {
	return repositories.PaymentHistoryRepository.FindPageByParams(db.GetDB(), params)
}

func (this *paymentHistoryService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.PaymentHistory, paging *commons.Paging) {
	return repositories.PaymentHistoryRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *paymentHistoryService) Create(t *models.PaymentHistory) error {
	return repositories.PaymentHistoryRepository.Create(db.GetDB(), t)
}

func (this *paymentHistoryService) Update(t *models.PaymentHistory) error {
	return repositories.PaymentHistoryRepository.Update(db.GetDB(), t)
}

func (this *paymentHistoryService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.PaymentHistoryRepository.Updates(db.GetDB(), id, columns)
}

func (this *paymentHistoryService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.PaymentHistoryRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *paymentHistoryService) Delete(id int64) {
	repositories.PaymentHistoryRepository.Delete(db.GetDB(), id)
}
