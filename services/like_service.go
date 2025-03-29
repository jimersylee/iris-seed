
package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var LikeService = newLikeService()

func newLikeService() *likeService {
	return &likeService {}
}

type likeService struct {
}

func (this *likeService) Get(id int64) *models.Like {
	return repositories.LikeRepository.Get(db.GetDB(), id)
}

func (this *likeService) Take(where ...interface{}) *models.Like {
	return repositories.LikeRepository.Take(db.GetDB(), where...)
}

func (this *likeService) Find(cnd *commons.SqlCnd) []models.Like {
	return repositories.LikeRepository.Find(db.GetDB(), cnd)
}

func (this *likeService) FindOne(cnd *commons.SqlCnd) *models.Like {
	return repositories.LikeRepository.FindOne(db.GetDB(), cnd)
}

func (this *likeService) FindPageByParams(params *commons.QueryParams) (list []models.Like, paging *commons.Paging) {
	return repositories.LikeRepository.FindPageByParams(db.GetDB(), params)
}

func (this *likeService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Like, paging *commons.Paging) {
	return repositories.LikeRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *likeService) Create(t *models.Like) error {
	return repositories.LikeRepository.Create(db.GetDB(), t)
}

func (this *likeService) Update(t *models.Like) error {
	return repositories.LikeRepository.Update(db.GetDB(), t)
}

func (this *likeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.LikeRepository.Updates(db.GetDB(), id, columns)
}

func (this *likeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.LikeRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *likeService) Delete(id int64) {
	repositories.LikeRepository.Delete(db.GetDB(), id)
}

