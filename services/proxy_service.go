package services

import (
	"crypto/tls"
	"errors"
	"github.com/jimersylee/iris-seed/commons"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net"
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
	//todo：找出能用的ip
	ipModel := IpService.FindOne(commons.NewSqlCnd().Where("status=1").Asc("request_times"))
	logrus.Infof("ipModel:%s,uri:%s,uriNeed:%s,webUrl:%s", ipModel, uri, uriNeed, webUrl)
	if ipModel == nil {
		logrus.Error("can't find available ip")
		ctx.ResponseWriter().WriteHeader(500)
		return
	}
	ip := ipModel.Ip
	IpService.incrRequestTimes(ip)
	responseStr, statusCode := p.fly(ip, webUrl)
	logrus.Infof("steam返回数据：%s", responseStr)
	ctx.ResponseWriter().WriteHeader(statusCode)
	_, _ = ctx.WriteString(responseStr)

}

func (p *ProxyServiceImpl) fly(ip string, webUrl string) (content string, statusCode int) {
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
		Timeout:   time.Second * 10, //超时时间
	}

	resp, err := client.Do(request)
	if err != nil {
		logrus.Errorf("访问steam出错，error:%s", err)
		return "", 600
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	return bodyString, resp.StatusCode
}

//let agent change ip
func (this *ProxyServiceImpl) ChangeIp(ip string) {
	// 调用管理接口换ip
	changeIpUrl := "http://" + ip + ":60001/adsl-start"
	get, err := this._get(changeIpUrl, nil, nil)
	if err != nil {
		logrus.Errorf("ChangeIp err:%s", err)
	}
	logrus.Infof("ChangeIp response:%s", get.Body)

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

//检测ip连通状态
func (this *ProxyServiceImpl) CheckIpAlive() {

	results := IpService.Find(commons.NewSqlCnd())
	for _, v := range results {
		address := v.Ip + ":" + strconv.Itoa(v.Port)
		dial, err := net.Dial("tcp", address)
		if err != nil {
			logrus.Warnf("this ip down：%s,err:%s", address, err)
			//既然不通了，那就删了
			IpService.Delete(v.ID)
			continue
		}
		defer dial.Close()

	}

}

//检查ip状态，429，500等统计数据
func (this *ProxyServiceImpl) CheckIpStatus() {

}

//删除无用的ip
func (this *ProxyServiceImpl) DeleteUselessIp() {

}
