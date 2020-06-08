package services

import (
	"crypto/tls"
	"github.com/jimersylee/iris-seed/commons"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
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
	responseStr, statusCode := p.fly(ip, webUrl)
	logrus.Infof("steam返回数据：%s", responseStr)
	ctx.ResponseWriter().WriteHeader(statusCode)
	ctx.WriteString(responseStr)

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
	return bodyString, resp.StatusCode
}
