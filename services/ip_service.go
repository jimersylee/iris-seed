package services

import (
	"crypto/tls"
	"fmt"
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/models"
	"github.com/jimersylee/iris-seed/repositories"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
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
	temp := repositories.IpRepository.FindOne(db.GetDB(), commons.NewSqlCnd().Eq("ip", "127.0.0.2"))
	logrus.Info("find result", temp)
	if temp != nil {
		return nil
	}
	t.CreateAt = time.Now()
	t.UpdateAt = time.Now()
	return repositories.IpRepository.Create(db.GetDB(), t)
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

func (i *ipService) Proxy(ip string, webUrl string) (content string, statusCode int) {
	proxyUrl := "http://" + ip + ":60002"
	logrus.Info("use " + proxyUrl + " to proxy")
	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	request, _ := http.NewRequest("GET", webUrl, nil)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Do(request)
	if err != nil {
		logrus.Errorf("访问steam出错，error:%s", err)
		return "", 500
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	fmt.Println(bodyString)
	return bodyString, resp.StatusCode
}
