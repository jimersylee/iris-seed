package services

import (
	"crypto/tls"
	"errors"
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/services/cache"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ProxyInterface interface {
	//注册ip
	ProxyCommunity(ctx iris.Context) bool
	ProxyApi(ctx iris.Context)
}

var ProxyService = NewProxyService()

func NewProxyService() *ProxyServiceImpl {
	return &ProxyServiceImpl{}
}

type ProxyServiceImpl struct {
}

func (p *ProxyServiceImpl) Proxy(ctx iris.Context) {
	uri := ctx.Request().RequestURI

	//判断是代理api还是社区
	uriNeed := ""
	webUrl := ""
	if strings.Index(uri, "steamapi") > 0 {
		uriNeed = uri[len("/api/steamapi/"):]
		webUrl = "http://api.steampowered.com/" + uriNeed
	} else {
		uriNeed = uri[len("/api/steamcommunity/"):]
		webUrl = "http://steamcommunity.com/" + uriNeed
	}
	//找出能用的ip
	ipModel := IpService.FindOne(commons.NewSqlCnd().Where("status=1").Asc("request_times"))
	logrus.Infof("ipModel:%s,uri:%s,uriNeed:%s,webUrl:%s", ipModel, uri, uriNeed, webUrl)
	if ipModel == nil {
		logrus.Error("can't find available ip")
		ctx.ResponseWriter().WriteHeader(500)
		return
	}
	ip := ipModel.Ip
	IpService.incrRequestTimes(ip)
	res := p.fly(ip, webUrl)
	cache.ProxyCache.IncrHttpStatusTimesByIpAndStatus(ip, res.StatusCode)

	for key, value := range res.Header {
		for _, v := range value {
			ctx.ResponseWriter().Header().Add(key, v)
		}
	}

	ctx.ResponseWriter().WriteHeader(res.StatusCode)

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		all, _ := ioutil.ReadAll(res.Body)
		body := string(all)
		logrus.Info("body=======" + body)
		_, _ = ctx.ResponseWriter().WriteString(body)
	} else {
		io.Copy(ctx.ResponseWriter(), res.Body)
	}
	res.Body.Close()

}

func (p *ProxyServiceImpl) fly(ip string, webUrl string) *http.Response {
	proxyUrl := "http://" + ip + ":60002"
	logrus.Info("use " + proxyUrl + " to proxy")
	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	request, _ := http.NewRequest("GET", webUrl, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10, //超时时间
	}

	resp, err := client.Do(request)
	if err != nil {
		logrus.Errorf("访问steam出错，error:%s", err)
		return nil
	}
	defer client.CloseIdleConnections()
	return resp
}

//let agent change ip
func (this *ProxyServiceImpl) ChangeIp(ip string) {
	// 调用管理接口换ip
	changeIpUrl := "http://" + ip + ":60001/adsl-start"
	get, err := this._get(changeIpUrl, nil, nil)
	if err != nil {
		logrus.Errorf("ChangeIp err:[%s]", err)
		return
	}
	logrus.Infof("ChangeIp response:[%s]", get.Body)

}
func (this *ProxyServiceImpl) _get(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	//new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail ")
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{}
	logrus.Infof("Go %s URL : %s \n", http.MethodGet, req.URL.String())
	return client.Do(req)
}

func (this *ProxyServiceImpl) AllCheckTask() {
	this.CheckIpAlive()
	this.CheckIpStatus()
	this.DeleteUselessIp()
}

//检测ip连通状态
func (this *ProxyServiceImpl) CheckIpAlive() {
	logrus.Info("start to check ip Alive")
	results := IpService.Find(commons.NewSqlCnd())
	for _, v := range results {
		address := v.Ip + ":" + strconv.Itoa(v.Port)
		logrus.Debug("check ip alive,address:" + address)
		urlToCheck := "http://" + address
		_, err := this._get(urlToCheck, nil, nil)
		if err != nil {
			logrus.Warnf("this ip down：%s,err:%s", address, err)
			//既然不通了，那就删了
			IpService.Delete(v.ID)
			continue
		}

	}
	logrus.Info("finish to check ip Alive")

}

//检查ip状态，429，500等统计数据
func (this *ProxyServiceImpl) CheckIpStatus() {
	logrus.Info("start to check ip status")
	all := cache.ProxyCache.IpPoolGetAll()
	logrus.Info("CheckIpStatus:", all)
	for _, ip := range all {
		needToBeBanned := cache.ProxyCache.CalcIpNeedToBeBanned(ip)
		if needToBeBanned == true {
			//从数据库中删除
			IpService.DeleteByIp(ip)
			this.ChangeIp(ip)
		}
	}
}

//定时删除无用的ip，N小时都没使用过的
func (this *ProxyServiceImpl) DeleteUselessIp() {
	logrus.Info("start to delete useless ip")
	ips := IpService.Find(commons.NewSqlCnd().Lt("update_at", time.Now().Unix()-3600*1))
	for _, v := range ips {
		IpService.Delete(v.ID)
	}
}
