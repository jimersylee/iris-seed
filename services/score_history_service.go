package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
)

var ScoreHistoryService = newScoreHistoryService()

func newScoreHistoryService() *scoreHistoryService {
	return &scoreHistoryService{}
}

type scoreHistoryService struct {
}

func (this *scoreHistoryService) Get(id int64) *models.ScoreHistory {
	return repositories.ScoreHistoryRepository.Get(db.GetDB(), id)
}

func (this *scoreHistoryService) Take(where ...interface{}) *models.ScoreHistory {
	return repositories.ScoreHistoryRepository.Take(db.GetDB(), where...)
}

func (this *scoreHistoryService) Find(cnd *commons.SqlCnd) []models.ScoreHistory {
	return repositories.ScoreHistoryRepository.Find(db.GetDB(), cnd)
}

func (this *scoreHistoryService) FindOne(cnd *commons.SqlCnd) *models.ScoreHistory {
	return repositories.ScoreHistoryRepository.FindOne(db.GetDB(), cnd)
}

func (this *scoreHistoryService) FindPageByParams(params *commons.QueryParams) (list []models.ScoreHistory, paging *commons.Paging) {
	return repositories.ScoreHistoryRepository.FindPageByParams(db.GetDB(), params)
}

func (this *scoreHistoryService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.ScoreHistory, paging *commons.Paging) {
	return repositories.ScoreHistoryRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *scoreHistoryService) Create(t *models.ScoreHistory) error {
	return repositories.ScoreHistoryRepository.Create(db.GetDB(), t)
}

func (this *scoreHistoryService) Update(t *models.ScoreHistory) error {
	return repositories.ScoreHistoryRepository.Update(db.GetDB(), t)
}

func (this *scoreHistoryService) Delete(id int64) {
	repositories.ScoreHistoryRepository.Delete(db.GetDB(), id)
}
