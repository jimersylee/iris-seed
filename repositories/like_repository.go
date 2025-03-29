
package repositories

import (
	"github.com/jimersylee/iris-seed/commons"	
	"github.com/jimersylee/iris-seed/models"
	"github.com/jinzhu/gorm"
)

var LikeRepository = newLikeRepository()

func newLikeRepository() *likeRepository {
	return &likeRepository{}
}

type likeRepository struct {
}

func (this *likeRepository) Get(db *gorm.DB, id int64) *models.Like {
	ret := &models.Like{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *likeRepository) Take(db *gorm.DB, where ...interface{}) *models.Like {
	ret := &models.Like{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *likeRepository) Find(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Like) {
	cnd.Find(db, &list)
	return
}

func (this *likeRepository) FindOne(db *gorm.DB, cnd *commons.SqlCnd) ( ret *models.Like) {
	cnd.FindOne(db, &ret)
	return 
}

func (this *likeRepository) FindPageByParams(db *gorm.DB, params *commons.QueryParams) (list []models.Like, paging *commons.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *likeRepository) FindPageByCnd(db *gorm.DB, cnd *commons.SqlCnd) (list []models.Like, paging *commons.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Like{})

	paging = &commons.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *likeRepository) Create(db *gorm.DB, t *models.Like) (err error) {
	err = db.Create(t).Error
	return
}

func (this *likeRepository) Update(db *gorm.DB, t *models.Like) (err error) {
	err = db.Save(t).Error
	return
}

func (this *likeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Like{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *likeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Like{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *likeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Like{}, "id = ?", id)
}

