package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var BookService = newBookService()

func newBookService() *bookService {
	return &bookService{}
}

type bookService struct {
}

func (this *bookService) Get(id int64) *models.Book {
	return repositories.BookRepository.Get(db.GetDB(), id)
}

func (this *bookService) Take(where ...interface{}) *models.Book {
	return repositories.BookRepository.Take(db.GetDB(), where...)
}

func (this *bookService) Find(cnd *commons.SqlCnd) []models.Book {
	return repositories.BookRepository.Find(db.GetDB(), cnd)
}

func (this *bookService) FindOne(cnd *commons.SqlCnd) *models.Book {
	return repositories.BookRepository.FindOne(db.GetDB(), cnd)
}

func (this *bookService) FindPageByParams(params *commons.QueryParams) (list []models.Book, paging *commons.Paging) {
	return repositories.BookRepository.FindPageByParams(db.GetDB(), params)
}

func (this *bookService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Book, paging *commons.Paging) {
	return repositories.BookRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *bookService) Create(t *models.Book) error {
	return repositories.BookRepository.Create(db.GetDB(), t)
}

func (this *bookService) Update(t *models.Book) error {
	return repositories.BookRepository.Update(db.GetDB(), t)
}

func (this *bookService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.BookRepository.Updates(db.GetDB(), id, columns)
}

func (this *bookService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.BookRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *bookService) Delete(id int64) {
	repositories.BookRepository.Delete(db.GetDB(), id)
}
