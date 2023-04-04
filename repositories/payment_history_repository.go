package repositories

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var PaymentHistoryRepository = newPaymentHistoryRepository()

func newPaymentHistoryRepository() *paymentHistoryRepository {
	return &paymentHistoryRepository{}
}

type paymentHistoryRepository struct {
}

func (this *paymentHistoryRepository) Get(db *gorm.DB, id int64) *models.PaymentHistory {
	ret := &models.PaymentHistory{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *paymentHistoryRepository) Take(db *gorm.DB, where ...interface{}) *models.PaymentHistory {
	ret := &models.PaymentHistory{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *paymentHistoryRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.PaymentHistory) {
	cnd.Find(db, &list)
	return
}

func (this *paymentHistoryRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) (ret *models.PaymentHistory) {
	cnd.FindOne(db, &ret)
	return
}

func (this *paymentHistoryRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.PaymentHistory, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *paymentHistoryRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.PaymentHistory, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.PaymentHistory{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *paymentHistoryRepository) Create(db *gorm.DB, t *models.PaymentHistory) (err error) {
	err = db.Create(t).Error
	return
}

func (this *paymentHistoryRepository) Update(db *gorm.DB, t *models.PaymentHistory) (err error) {
	err = db.Save(t).Error
	return
}

func (this *paymentHistoryRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.PaymentHistory{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *paymentHistoryRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.PaymentHistory{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *paymentHistoryRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.PaymentHistory{}, "id = ?", id)
}
