package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var CommentService = newCommentService()

func newCommentService() *commentService {
	return &commentService{}
}

type commentService struct {
}

func (this *commentService) Get(id int64) *models.Comment {
	return repositories.CommentRepository.Get(db.GetDB(), id)
}

func (this *commentService) Take(where ...interface{}) *models.Comment {
	return repositories.CommentRepository.Take(db.GetDB(), where...)
}

func (this *commentService) Find(cnd *commons.SqlCnd) []models.Comment {
	return repositories.CommentRepository.Find(db.GetDB(), cnd)
}

func (this *commentService) FindOne(cnd *commons.SqlCnd) *models.Comment {
	return repositories.CommentRepository.FindOne(db.GetDB(), cnd)
}

func (this *commentService) FindPageByParams(params *commons.QueryParams) (list []models.Comment, paging *commons.Paging) {
	return repositories.CommentRepository.FindPageByParams(db.GetDB(), params)
}

func (this *commentService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Comment, paging *commons.Paging) {
	return repositories.CommentRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *commentService) Create(t *models.Comment) error {
	return repositories.CommentRepository.Create(db.GetDB(), t)
}

func (this *commentService) Update(t *models.Comment) error {
	return repositories.CommentRepository.Update(db.GetDB(), t)
}

func (this *commentService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CommentRepository.Updates(db.GetDB(), id, columns)
}

func (this *commentService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CommentRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *commentService) Delete(id int64) {
	repositories.CommentRepository.Delete(db.GetDB(), id)
}
