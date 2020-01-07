package repositories

import (
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jinzhu/gorm"
)

func NewUserTokenRepository() *userTokenRepository {
	return &userTokenRepository{}
}

type userTokenRepository struct {
}
type UserTokenRepository interface {
	InsertOne(db *gorm.DB, userToken datamodels.UserToken) (err error)
	Delete(db *gorm.DB, id int64)
	FindByToken(db *gorm.DB, token string) *datamodels.UserToken
	FindOne(db *gorm.DB, id int64) (user *datamodels.UserToken)
}

func (this *userTokenRepository) InsertOne(db *gorm.DB, userToken datamodels.UserToken) (err error) {
	err = db.Create(userToken).Error
	return
}

func (this *userTokenRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&datamodels.UserToken{}, "id = ?", id)
}

func (this *userTokenRepository) FindByToken(db *gorm.DB, token string) *datamodels.UserToken {
	if len(token) == 0 {
		return nil
	}
	return this.Take(db, "token = ?", token)
}

func (this *userTokenRepository) FindOne(db *gorm.DB, id int64) *datamodels.UserToken {
	ret := &datamodels.UserToken{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userTokenRepository) Take(db *gorm.DB, where ...interface{}) *datamodels.UserToken {
	ret := &datamodels.UserToken{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}
