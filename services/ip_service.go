package services

import (
	"github.com/jimersylee/go-steam-proxy/commons"
	"github.com/jimersylee/go-steam-proxy/commons/db"
	"github.com/jimersylee/go-steam-proxy/models"
	"github.com/jimersylee/go-steam-proxy/repositories"
	"github.com/jimersylee/go-steam-proxy/services/cache"
	"github.com/sirupsen/logrus"
	"time"
)

var IpService = newIpService()

func newIpService() *ipService {
	return &ipService{}
}

type ipService struct {
}

func (this *ipService) Get(id int64) *models.Ip {
	return repositories.IpRepository.Get(db.GetDB(), id)
}

func (this *ipService) Take(where ...interface{}) *models.Ip {
	return repositories.IpRepository.Take(db.GetDB(), where...)
}

func (this *ipService) Find(cnd *commons.SqlCnd) []models.Ip {
	return repositories.IpRepository.Find(db.GetDB(), cnd)
}

func (this *ipService) FindOne(cnd *commons.SqlCnd) *models.Ip {
	return repositories.IpRepository.FindOne(db.GetDB(), cnd)
}

func (this *ipService) FindPageByParams(params *commons.QueryParams) (list []models.Ip, paging *commons.Paging) {
	return repositories.IpRepository.FindPageByParams(db.GetDB(), params)
}

func (this *ipService) FindPageByCnd(cnd *commons.SqlCnd) (list []models.Ip, paging *commons.Paging) {
	return repositories.IpRepository.FindPageByCnd(db.GetDB(), cnd)
}

func (this *ipService) Create(t *models.Ip) error {
	cache.ProxyCache.IpPoolAdd(t.Ip)
	temp := repositories.IpRepository.FindOne(db.GetDB(), commons.NewSqlCnd().Eq("ip", t.Ip))
	logrus.Info("find result", temp)
	if temp != nil {
		return nil
	}
	t.CreateAt = time.Now().Unix()
	t.UpdateAt = time.Now().Unix()
	t.Status = 1
	t.Port = 60002
	err := repositories.IpRepository.Create(db.GetDB(), t)
	if err != nil {
		return err
	}

	return nil
}

func (this *ipService) Update(t *models.Ip) error {
	return repositories.IpRepository.Update(db.GetDB(), t)
}

func (this *ipService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.IpRepository.Updates(db.GetDB(), id, columns)
}

func (this *ipService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.IpRepository.UpdateColumn(db.GetDB(), id, name, value)
}

func (this *ipService) Delete(id int64) {
	repositories.IpRepository.Delete(db.GetDB(), id)
}
func (this *ipService) DeleteByIp(ip string) {
	one := repositories.IpRepository.FindOne(db.GetDB(), commons.NewSqlCnd().Eq("ip", ip))
	if one != nil {
		repositories.IpRepository.Delete(db.GetDB(), one.ID)
	}

}

func (this *ipService) incrRequestTimes(ip string) {
	modelIp := repositories.IpRepository.FindOne(db.GetDB(), commons.NewSqlCnd().Eq("ip", ip))
	if modelIp == nil {
		return
	}
	modelIp.RequestTimes = modelIp.RequestTimes + 1
	modelIp.UpdateAt = time.Now().Unix()
	err := repositories.IpRepository.Update(db.GetDB(), modelIp)
	if err != nil {
		logrus.Errorf("incrRequestTimes failed,err:%s", err)
	}
}
