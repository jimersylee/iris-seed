
package repositories

import (
	"github.com/jimersylee/iris-seed/commons"	
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var ScoreHistoryRepository = newScoreHistoryRepository()

func newScoreHistoryRepository() *scoreHistoryRepository {
	return &scoreHistoryRepository{}
}

type scoreHistoryRepository struct {
}

func (this *scoreHistoryRepository) Get(db *gorm.DB, id int64) *models.ScoreHistory {
	ret := &models.ScoreHistory{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *scoreHistoryRepository) Take(db *gorm.DB, where ...interface{}) *models.ScoreHistory {
	ret := &models.ScoreHistory{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *scoreHistoryRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.ScoreHistory) {
	cnd.Find(db, &list)
	return
}

func (this *scoreHistoryRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) ( ret *models.ScoreHistory) {
	cnd.FindOne(db, &ret)
	return 
}

func (this *scoreHistoryRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.ScoreHistory, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *scoreHistoryRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.ScoreHistory, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.ScoreHistory{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *scoreHistoryRepository) Create(db *gorm.DB, t *models.ScoreHistory) (err error) {
	err = db.Create(t).Error
	return
}

func (this *scoreHistoryRepository) Update(db *gorm.DB, t *models.ScoreHistory) (err error) {
	err = db.Save(t).Error
	return
}

func (this *scoreHistoryRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.ScoreHistory{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *scoreHistoryRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.ScoreHistory{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *scoreHistoryRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.ScoreHistory{}, "id = ?", id)
}

