package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var NoteService = newNoteService()

func newNoteService() *noteService {
	return &noteService{}
}

type noteService struct {
}

func (this *noteService) Get(id int64) *models.Note {
	return repositories.NoteRepository.Get(db.GetDB(), id)
}

func (this *noteService) Take(where ...interface{}) *models.Note {
	return repositories.NoteRepository.Take(db.GetDB(), where...)
}

func (this *noteService) Find(cnd *commons.SqlCnd) []models.Note {
	return repositories.NoteRepository.Find(db.GetDB(), cnd)
}

func (this *noteService) FindOne(cnd *commons.SqlCnd) *models.Note {
	return repositories.NoteRepository.FindOne(db.GetDB(), cnd)
}

func (this *noteService) FindPageByParams(params *commons.QueryParams) (list []models.Note, paging *commons.Paging) {
	return repositories.NoteRepository.FindPageByParams(db.GetDB(), params)
}

func (this *noteService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Note, paging *commons.Paging) {
	return repositories.NoteRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *noteService) Create(t *models.Note) error {
	return repositories.NoteRepository.Create(db.GetDB(), t)
}

func (this *noteService) Update(t *models.Note) error {
	return repositories.NoteRepository.Update(db.GetDB(), t)
}

func (this *noteService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.NoteRepository.Updates(db.GetDB(), id, columns)
}

func (this *noteService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.NoteRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *noteService) Delete(id int64) {
	repositories.NoteRepository.Delete(db.GetDB(), id)
}
