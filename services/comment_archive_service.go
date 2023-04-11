package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var CommentArchiveService = newCommentArchiveService()

func newCommentArchiveService() *commentArchiveService {
	return &commentArchiveService{}
}

type commentArchiveService struct {
}

func (this *commentArchiveService) Get(id int64) *models.CommentArchive {
	return repositories.CommentArchiveRepository.Get(db.GetDB(), id)
}

func (this *commentArchiveService) Take(where ...interface{}) *models.CommentArchive {
	return repositories.CommentArchiveRepository.Take(db.GetDB(), where...)
}

func (this *commentArchiveService) Find(cnd *commons.SqlCnd) []models.CommentArchive {
	return repositories.CommentArchiveRepository.Find(db.GetDB(), cnd)
}

func (this *commentArchiveService) FindOne(cnd *commons.SqlCnd) *models.CommentArchive {
	return repositories.CommentArchiveRepository.FindOne(db.GetDB(), cnd)
}

func (this *commentArchiveService) FindPageByParams(params *commons.QueryParams) (list []models.CommentArchive, paging *commons.Paging) {
	return repositories.CommentArchiveRepository.FindPageByParams(db.GetDB(), params)
}

func (this *commentArchiveService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.CommentArchive, paging *commons.Paging) {
	return repositories.CommentArchiveRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *commentArchiveService) Create(t *models.CommentArchive) error {
	return repositories.CommentArchiveRepository.Create(db.GetDB(), t)
}

func (this *commentArchiveService) Update(t *models.CommentArchive) error {
	return repositories.CommentArchiveRepository.Update(db.GetDB(), t)
}

func (this *commentArchiveService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CommentArchiveRepository.Updates(db.GetDB(), id, columns)
}

func (this *commentArchiveService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CommentArchiveRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *commentArchiveService) Delete(id int64) {
	repositories.CommentArchiveRepository.Delete(db.GetDB(), id)
}
